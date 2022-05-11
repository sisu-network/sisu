package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerSlashValidator struct {
	pmm        PostedMessageManager
	mc         ManagerContainer
	keeper     keeper.Keeper
	valManager ValidatorManager
}

func NewHandlerSlashValidator(mc ManagerContainer) *HandlerSlashValidator {
	return &HandlerSlashValidator{
		pmm:        mc.PostedMessageManager(),
		mc:         mc,
		keeper:     mc.Keeper(),
		valManager: NewValidatorManager(mc.Keeper()),
	}
}

func (h *HandlerSlashValidator) DeliverMsg(ctx sdk.Context, msg *types.SlashValidatorMsg) (*sdk.Result, error) {
	vals := h.valManager.GetNodesByStatus(types.NodeStatus_Validator)
	for addr, _ := range vals {
		log.Debug("val address = ", addr)
	}

	nodeAddr, err := sdk.AccAddressFromBech32(msg.Data.NodeAddress)
	if err != nil {
		log.Error("error when parsing address ", msg.Data.NodeAddress)
		return &sdk.Result{}, err
	}

	if err := h.keeper.IncSlashToken(ctx, nodeAddr, msg.Data.SlashPoint); err != nil {
		log.Error(err)
		return &sdk.Result{}, err
	}

	afterSlashBalance, err := h.keeper.GetSlashToken(ctx, nodeAddr)
	if err != nil {
		log.Error(err)
		return &sdk.Result{}, err
	}

	log.Debugf("after slash balance of node %s is %d", nodeAddr, afterSlashBalance)
	return &sdk.Result{}, nil
}
