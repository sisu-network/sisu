package components

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

type PostedMessageManager interface {
	ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte)
}

type defaultPostedMessageManager struct {
	keeper keeper.Keeper
}

func NewPostedMessageManager(keeper keeper.Keeper) *defaultPostedMessageManager {
	return &defaultPostedMessageManager{
		keeper: keeper,
	}
}

func (m *defaultPostedMessageManager) ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte) {
	hash, signer, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Error("failed to get tx hash, err = ", err)
		return false, hash
	}

	if m.keeper.IsTxRecordProcessed(ctx, hash) {
		return false, hash
	}

	count := m.keeper.SaveTxRecord(ctx, hash, signer)
	tssParams := m.keeper.GetParams(ctx)
	if tssParams == nil {
		log.Warn("tssParams is nil")
		return false, nil
	}

	if count >= int(tssParams.MajorityThreshold) {
		return true, hash
	}

	return false, hash
}
