package sisu

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
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
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doReshareResult(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerReshareResult) doReshareResult(ctx sdk.Context, msg *types.ReshareResultWithSigner) ([]byte, error) {
	newValPubKeys := msg.Data.NewValidatorSetPubKeyBytes
	// TODO: check validator address who sent confirm msg is in new validator set
	if isConfirmedByAllVals, _ := h.pmm.IsReachedThreshold(
		ctx,
		msg,
		len(newValPubKeys),
	); !isConfirmedByAllVals {
		return nil, nil
	}

	nodes := make([]*types.Node, len(newValPubKeys))
	for _, bz := range newValPubKeys {
		cosmosPubKey, err := utils.GetCosmosPubKey("ed25519", bz)
		if err != nil {
			log.Error("error when get cosmos public key. error = ", err)
			return nil, err
		}

		node := &types.Node{
			Id: hex.EncodeToString(cosmosPubKey.Address()),
			ConsensusKey: &types.Pubkey{
				Type:  cosmosPubKey.Type(),
				Bytes: cosmosPubKey.Bytes(),
			},
			AccAddress:  sdk.AccAddress(cosmosPubKey.Address()).String(),
			IsValidator: true,
			Status:      types.NodeStatus_Validator,
		}

		nodes = append(nodes, node)
	}

	if err := h.valManager.SetValidators(ctx, nodes); err != nil {
		return nil, err
	}

	return []byte{}, nil
}
