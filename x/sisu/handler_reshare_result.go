package sisu

import (
	"bytes"

	cmcrypto "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"golang.org/x/exp/slices"
)

type HandlerReshareResult struct {
	mc         ManagerContainer
	pmm        PostedMessageManager
	keeper     keeper.Keeper
	valManager ValidatorManager
}

func NewHandlerReshareResult(mc ManagerContainer) *HandlerReshareResult {
	return &HandlerReshareResult{
		mc:         mc,
		pmm:        mc.PostedMessageManager(),
		keeper:     mc.Keeper(),
		valManager: mc.ValidatorManager(),
	}
}

func (h *HandlerReshareResult) DeliverMsg(ctx sdk.Context, msg *types.ReshareResultWithSigner) (*sdk.Result, error) {
	log.Debug("handling HandlerReshareResult ...")

	rcHash, signer, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Error("error when getting tx record hash: ", err)
		return &sdk.Result{}, nil
	}

	h.keeper.SaveTxRecord(ctx, rcHash, signer)

	if h.keeper.IsTxRecordProcessed(ctx, rcHash) {
		return &sdk.Result{}, nil
	}

	if err := h.keeper.IncSlashToken(ctx, getSlashToken(msg), msg.GetSender()); err != nil {
		return &sdk.Result{}, nil
	}

	data, err := h.doReshareResult(ctx, msg, rcHash)
	if err != nil {
		return &sdk.Result{}, err
	}

	voters := h.keeper.GetVotersInAccAddress(ctx, rcHash)
	if err := h.keeper.DecSlashToken(ctx, getSlashToken(msg), voters...); err != nil {
		return &sdk.Result{}, err
	}

	return &sdk.Result{Data: data}, nil
}

// reshare result must be confirmed by all new validators
func (h *HandlerReshareResult) doReshareResult(ctx sdk.Context, msg *types.ReshareResultWithSigner, rcHash []byte) ([]byte, error) {
	newValPubKeys := msg.Data.NewValidatorSetPubKeyBytes
	// Get all voters who signed this tx
	voters := h.keeper.GetVoters(ctx, rcHash)

	// newVoters is set of new validator's AccAddress
	newVoters := make([]string, 0)
	nodes := h.valManager.GetNodesByStatus(types.NodeStatus_Unknown)
	if len(nodes) == 0 {
		return nil, nil
	}

	// converts from consensus key to acc address
	for _, valPk := range newValPubKeys {
		node := nodes[string(valPk)]
		newVoters = append(newVoters, node.AccAddress)
	}

	for _, v := range newVoters {
		if !slices.Contains(voters, v) {
			log.Debug("doReshareResult: reshare result have not enough vote from new validator set. Ignore msg")
			return nil, nil
		}
	}

	log.Debug("doReshareResult: enough vote. Process reshare result")

	// Build []abci.ValidatorUpdate
	// Dheart is using cosmos pubkey format as input, meanwhile abci.ValidatorUpdate is using Tendermint pubkey format
	valUpdates := make([]abci.ValidatorUpdate, 0)
	for _, pk := range newValPubKeys {
		cosmosPubKey, err := utils.GetCosmosPubKey("ed25519", pk)
		if err != nil {
			log.Error("error when get cosmos public key. error = ", err)
			return nil, err
		}

		tmProtoPk, err := cmcrypto.ToTmProtoPublicKey(cosmosPubKey)
		if err != nil {
			log.Error("error when convert from cosmos pubkey to tendermint pubkey. error = ", err)
			return nil, err
		}

		message := abci.Ed25519ValidatorUpdate(tmProtoPk.GetEd25519(), 100)
		valUpdates = append(valUpdates, message)
	}

	// Reset old validators
	currentValidators := h.keeper.LoadNodesByStatus(ctx, types.NodeStatus_Validator)
	for _, v := range currentValidators {
		isNewValidator := false
		for _, pk := range newValPubKeys {
			if bytes.Equal(v.ConsensusKey.GetBytes(), pk) {
				isNewValidator = true
				break
			}
		}

		if isNewValidator {
			continue
		}

		cosmosPubKey, err := utils.GetCosmosPubKey("ed25519", v.ConsensusKey.GetBytes())
		if err != nil {
			log.Error("error when get cosmos public key. error = ", err)
			return nil, err
		}

		tmProtoPk, err := cmcrypto.ToTmProtoPublicKey(cosmosPubKey)
		if err != nil {
			log.Error("error when convert from cosmos pubkey to tendermint pubkey. error = ", err)
			return nil, err
		}

		message := abci.Ed25519ValidatorUpdate(tmProtoPk.GetEd25519(), 0)
		valUpdates = append(valUpdates, message)
	}

	if len(valUpdates) == 0 {
		log.Debug("valUpdates is empty")
		return []byte{}, nil
	}

	if err := h.keeper.SaveIncomingValidatorUpdates(ctx, valUpdates); err != nil {
		return nil, err
	}

	log.Debug("SaveIncomingValidatorUpdates successfully", valUpdates)

	afterSavedValUpdate := h.keeper.GetIncomingValidatorUpdates(ctx)
	log.Debug("len(afterSavedValUpdate) = ", len(afterSavedValUpdate))
	h.keeper.ProcessTxRecord(ctx, rcHash)

	return []byte{}, nil
}
