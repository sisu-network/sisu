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
		mc:     mc,
	}
}

func (h *HandlerKeygen) DeliverMsg(ctx sdk.Context, signerMsg *types.KeygenWithSigner) (*sdk.Result, error) {
	rcHash, _, err := keeper.GetTxRecordHash(signerMsg)
	if err != nil {
		return &sdk.Result{}, err
	}

	if h.keeper.IsTxRecordProcessed(ctx, rcHash) {
		return &sdk.Result{}, nil
	}

	log.Debug("IncSlashToken in HandlerKeygen")
	if err := h.keeper.IncSlashToken(ctx, types.ObserveSlashPoint, signerMsg.GetSender()); err != nil {
		return &sdk.Result{}, nil
	}

	a, err := h.keeper.GetSlashToken(ctx, signerMsg.GetSender())
	if err != nil {
		return &sdk.Result{}, err
	}

	log.Debug("After IncSlashToken in HandlerKeygen. slash = ", a)

	log.Info("Delivering keygen, signer = ", signerMsg.Signer)
	pmm := h.mc.PostedMessageManager()
	if process, hash := pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doKeygen(ctx, signerMsg)
		if err != nil {
			return &sdk.Result{}, err
		}

		h.keeper.ProcessTxRecord(ctx, hash)

		voters := h.keeper.GetVotersInAccAddress(ctx, hash)
		log.Debug("before dec slash token")
		for _, v := range voters {
			a, err := h.keeper.GetSlashToken(ctx, v)
			if err != nil {
				return &sdk.Result{}, err
			}

			log.Debug("before dec slash. v = ", v.String(), " slash = ", a)
		}

		log.Debug("voters = ", h.keeper.GetVoters(ctx, hash))
		if err := h.keeper.DecSlashToken(ctx, types.ObserveSlashPoint, voters...); err != nil {
			return &sdk.Result{}, err
		}

		log.Debug("after dec slash token")
		for _, v := range voters {
			a, err := h.keeper.GetSlashToken(ctx, v)
			if err != nil {
				return &sdk.Result{}, err
			}

			log.Debug("after dec slash. v = ", v.String(), " slash = ", a)
		}

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
