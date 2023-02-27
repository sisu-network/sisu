package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerKeygen struct {
	mc     background.ManagerContainer
	keeper keeper.Keeper
}

func NewHandlerKeygen(mc background.ManagerContainer) *HandlerKeygen {
	return &HandlerKeygen{
		keeper: mc.Keeper(),
		mc:     mc,
	}
}

func (h *HandlerKeygen) DeliverMsg(ctx sdk.Context, msg *types.KeygenWithSigner) (*sdk.Result, error) {
	log.Info("Delivering keygen, signer = ", msg.Signer, " type = ", msg.Data.KeyType)
	pmm := h.mc.PostedMessageManager()
	if process, hash := pmm.ShouldProcessMsg(ctx, msg); process {
		data, err := h.doKeygen(ctx, msg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerKeygen) doKeygen(ctx sdk.Context, signerMsg *types.KeygenWithSigner) ([]byte, error) {
	log.Info("Doing keygen....")

	msg := signerMsg.Data
	globalData := h.mc.GlobalData()

	// Save this into Keeper && private db.
	h.keeper.SaveKeygen(ctx, msg)

	if globalData.IsCatchingUp() {
		return nil, nil
	}

	// Invoke TSS keygen in dheart
	h.doTss(ctx, msg)

	return []byte{}, nil
}

func (h *HandlerKeygen) doTss(ctx sdk.Context, msg *types.Keygen) {
	log.Info("doing keygen tss...")

	valsMag := h.mc.ValidatorManager()
	dheartClient := h.mc.DheartClient()

	// Send a signal to Dheart to start keygen process.
	log.Info("Sending keygen request to Dheart. KeyType =", msg.KeyType)
	pubKeys := valsMag.GetValidatorPubkeys(ctx)
	keygenId := helper.GetKeygenId(msg.KeyType, ctx.BlockHeight(), pubKeys)

	err := dheartClient.KeyGen(keygenId, msg.KeyType, pubKeys)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Keygen request is sent successfully.")
}
