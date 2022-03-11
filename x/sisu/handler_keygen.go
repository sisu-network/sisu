package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerKeygen struct {
	mc     ManagerContainer
	keeper keeper.Keeper
}

func NewHandlerKeygen(mc ManagerContainer) *HandlerKeygen {
	return &HandlerKeygen{
		keeper: mc.Keeper(),
	}
}

func (h *HandlerKeygen) DeliverMsg(ctx sdk.Context, signerMsg *types.KeygenWithSigner) (*sdk.Result, error) {
	log.Info("Delivering keygen, signer = ", signerMsg.Signer)
	pmm := h.mc.PostedMessageManager()
	if process, hash := pmm.ShouldProcessMsg(ctx, signerMsg); process {
		h.doKeygen(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)
	}

	return nil, nil
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
	h.doTss(msg, ctx.BlockHeight())

	return []byte{}, nil
}

func (h *HandlerKeygen) doTss(msg *types.Keygen, blockHeight int64) {
	log.Info("doing keygen tsss...")

	partyManager := h.mc.PartyManager()
	dheartClient := h.mc.DheartClient()

	// Send a signal to Dheart to start keygen process.
	log.Info("Sending keygen request to Dheart. KeyType =", msg.KeyType)
	pubKeys := partyManager.GetActivePartyPubkeys()
	keygenId := helper.GetKeygenId(msg.KeyType, blockHeight, pubKeys)

	err := dheartClient.KeyGen(keygenId, msg.KeyType, pubKeys)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Keygen request is sent successfully.")
}
