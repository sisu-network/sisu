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
		valManager: NewValidatorManager(mc.Keeper()),
	}
}

func (h *HandlerReshareResult) DeliverMsg(ctx sdk.Context, signerMsg *types.ReshareResultWithSigner) (*sdk.Result, error) {
	log.Debug("handling HandlerReshareResult ...")
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doReshareResult(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerReshareResult) doReshareResult(ctx sdk.Context, msg *types.ReshareResultWithSigner) ([]byte, error) {
	newValPubKeys := msg.Data.NewValidatorSetPubKeyBytes
	rcHash, _, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Error("error when getting tx record hash: ", err)
		return nil, err
	}

	if vote := h.valManager.CountVote(ctx, rcHash); vote < len(msg.Data.NewValidatorSetPubKeyBytes) {
		log.Debug("reshare result is not confirmed by all new validator")
		return nil, nil
	}

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

	return []byte{}, nil
}
