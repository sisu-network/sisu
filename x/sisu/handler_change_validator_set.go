package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerChangeValidatorSet struct {
	pmm    PostedMessageManager
	mc     ManagerContainer
	keeper keeper.Keeper
}

func NewHandlerChangeValidatorSet(mc ManagerContainer) *HandlerChangeValidatorSet {
	return &HandlerChangeValidatorSet{
		pmm:    mc.PostedMessageManager(),
		mc:     mc,
		keeper: mc.Keeper(),
	}
}

func (h *HandlerChangeValidatorSet) DeliverMsg(ctx sdk.Context, msg *types.ChangeValidatorSetMsg) (*sdk.Result, error) {
	process, hash := h.pmm.ShouldProcessMsg(ctx, msg)
	if !process {
		return &sdk.Result{}, nil
	}

	if err := h.doChangeValidatorSet(msg); err != nil {
		return &sdk.Result{}, err
	}

	h.keeper.ProcessTxRecord(ctx, hash)
	return &sdk.Result{}, nil
}

func (h *HandlerChangeValidatorSet) doChangeValidatorSet(msg *types.ChangeValidatorSetMsg) error {
	if h.mc.GlobalData().IsCatchingUp() {
		log.Info("We are catching up with the network, exiting doChangeValidatorSet")
		return nil
	}

	oldPubKeys, newPubKeys, err := msg.GetOldAndNewValidatorSet()
	if err != nil {
		return err
	}

	dheartClient := h.mc.DheartClient()
	if err := dheartClient.Reshare(oldPubKeys, newPubKeys); err != nil {
		log.Error("error when sending reshape request to heart. error = ", err)
		return err
	}

	log.Info("Reshape request is sent successfully")
	return nil
}
