package sisu

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/cardano-be/src/handler"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// HandlerMintMultiAsset mints multi asset on Cardano chain
type HandlerMintMultiAsset struct {
	globalData common.GlobalData
	pmm        PostedMessageManager
	mc         ManagerContainer
	keeper     keeper.Keeper
}

func NewHandleMintMultiAsset(mc ManagerContainer) *HandlerMintMultiAsset {
	return &HandlerMintMultiAsset{
		pmm:        mc.PostedMessageManager(),
		keeper:     mc.Keeper(),
		globalData: mc.GlobalData(),
		mc:         mc,
	}
}

func (h *HandlerMintMultiAsset) DeliverMsg(ctx sdk.Context, msg *types.CardanoMintMultiAssetMsg) (*sdk.Result, error) {
	log.Debug("Delivering mint multi asset ...")
	if h.globalData.IsCatchingUp() {
		log.Info("We are catching up with the network, exiting mint multi asset")
		return nil, nil
	}

	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		data, err := h.doMintMultiAsset(msg)
		h.keeper.ProcessTxRecord(ctx, hash)
		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerMintMultiAsset) doMintMultiAsset(msg *types.CardanoMintMultiAssetMsg) ([]byte, error) {
	req := &handler.MintRequest{
		PubKeys:          []htypes.PubKeyWrapper{},
		Lovelace:         uint64(msg.Data.Lovelace),
		AssetName:        msg.Data.AssetName,
		TSSPubKey:        msg.Data.TssPubkey,
		MultiAssetAmount: msg.Data.AssetAmount,
		Network:          1, // 0: mainnet, 1: testnet
	}
	resp, err := h.mc.CardanoClient().MintMultiAsset(context.Background(), req)
	if err != nil {
		return nil, err
	}

	log.Debug("cardano client response = ", resp)
	return nil, nil
}
