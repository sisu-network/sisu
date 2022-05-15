package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerChangeValidatorSet struct {
	pmm        PostedMessageManager
	mc         ManagerContainer
	keeper     keeper.Keeper
	valManager ValidatorManager
	globalData common.GlobalData
}

func NewHandlerChangeValidatorSet(mc ManagerContainer) *HandlerChangeValidatorSet {
	return &HandlerChangeValidatorSet{
		pmm:        mc.PostedMessageManager(),
		mc:         mc,
		keeper:     mc.Keeper(),
		valManager: mc.ValidatorManager(),
		globalData: mc.GlobalData(),
	}
}

func (h *HandlerChangeValidatorSet) DeliverMsg(ctx sdk.Context, msg *types.ChangeValidatorSetMsg) (*sdk.Result, error) {
	shouldProcess, rcHash := h.pmm.ShouldProcessMsg(ctx, msg)
	if !shouldProcess {
		return &sdk.Result{}, nil
	}

	if err := h.doChangeValidatorSet(msg); err != nil {
		return &sdk.Result{}, err
	}

	h.keeper.ProcessTxRecord(ctx, rcHash)
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

	log.Debug("oldPubKeys = ", oldPubKeys)
	log.Debug("newPubKeys = ", newPubKeys)
	dheartClient := h.mc.DheartClient()
	if err := dheartClient.Reshare(oldPubKeys, newPubKeys); err != nil {
		log.Error("error when sending reshare request to heart. error = ", err)
		return err
	}

	log.Info("Reshare request is sent successfully")
	return nil
}
