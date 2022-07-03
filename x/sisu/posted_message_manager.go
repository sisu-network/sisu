package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

type PostedMessageManager interface {
	ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte)
}

type DefaultPostedMessageManager struct {
	keeper keeper.Keeper
}

func NewPostedMessageManager(keeper keeper.Keeper) *DefaultPostedMessageManager {
	return &DefaultPostedMessageManager{
		keeper: keeper,
	}
}

func (m *DefaultPostedMessageManager) ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte) {
	hash, signer, err := keeper.GetTxRecordHash(msg)
	fmt.Println("ShouldProcessMsg: err = ", err)

	if err != nil {
		log.Error("failed to get tx hash, err = ", err)
		return false, hash
	}

	count := m.keeper.SaveTxRecord(ctx, hash, signer)
	tssParams := m.keeper.GetParams(ctx)
	if tssParams == nil {
		log.Warn("tssParams is nil")
		return false, nil
	}

	fmt.Println("ShouldProcessMsg: tssParams.MajorityThreshold = ", tssParams.MajorityThreshold)
	fmt.Println("ShouldProcessMsg: m.keeper.IsTxRecordProcessed(ctx, hash) = ", m.keeper.IsTxRecordProcessed(ctx, hash))

	if count >= int(tssParams.MajorityThreshold) && !m.keeper.IsTxRecordProcessed(ctx, hash) {
		return true, hash
	}

	return false, hash
}
