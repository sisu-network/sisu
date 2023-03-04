package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOutProposal struct {
	pmm         components.PostedMessageManager
	keeper      keeper.Keeper
	valsManager components.ValidatorManager
	globalData  components.GlobalData
	txSubmit    components.TxSubmit
	appKeys     components.AppKeys
	privateDb   keeper.PrivateDb
	background  background.Background
}

func NewHandlerTxOutProposal(mc background.ManagerContainer) *HandlerTxOutProposal {
	return &HandlerTxOutProposal{
		keeper:      mc.Keeper(),
		pmm:         mc.PostedMessageManager(),
		valsManager: mc.ValidatorManager(),
		globalData:  mc.GlobalData(),
		txSubmit:    mc.TxSubmit(),
		appKeys:     mc.AppKeys(),
		privateDb:   mc.PrivateDb(),
		background:  mc.Background(),
	}
}

// There are 2 cases where a TxOut can be finalized:
// 1) The assigned validator submits the TxOut and it's approved 2/3 of validators
// 2) The proposed txOut is rejected or it is not produced during a timeout period. At this time,
// every validator node submits its own txOut and everyone to come up with a consensused txOut.
func (h *HandlerTxOutProposal) DeliverMsg(ctx sdk.Context, msg *types.TxOutMsg) (*sdk.Result, error) {
	txOut := msg.Data

	validatorId := txOut.GetValidatorId()
	if len(validatorId) == 0 {
		log.Errorf("Validator id is empty for txout")
		return &sdk.Result{}, nil
	}

	assignedNode, err := h.valsManager.GetAssignedValidator(ctx, validatorId)
	if err == nil && assignedNode.AccAddress == msg.Signer {
		// This is the proposed TxOut from the assigned validator.
		h.keeper.AddProposedTxOut(ctx, msg.Data)

		// Add this message to the validation queue.
		h.background.AddVoteTxOut(ctx.BlockHeight(), msg)
	}

	return &sdk.Result{}, nil
}
