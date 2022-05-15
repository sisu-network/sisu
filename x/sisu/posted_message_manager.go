package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

var _ PostedMessageManager = (*DefaultPostedMessageManager)(nil)

//go:generate mockgen -source=./x/sisu/posted_message_manager.go -destination=./tests/mock/x/sisu/posted_message_manager.go -package=mock
type PostedMessageManager interface {
	ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte)
	IsReachedThreshold(ctx sdk.Context, msg sdk.Msg, threshold int) (bool, []byte)
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
	tssParams := m.keeper.GetParams(ctx)
	if tssParams == nil {
		return false, nil
	}

	return m.IsReachedThreshold(ctx, msg, int(tssParams.MajorityThreshold))
}

func (m *DefaultPostedMessageManager) IsReachedThreshold(ctx sdk.Context, msg sdk.Msg, threshold int) (bool, []byte) {
	hash, signer, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Error("failed to get tx hash, err = ", err)
		return false, hash
	}

	count := m.keeper.SaveTxRecord(ctx, hash, signer)
	if count >= threshold && !m.keeper.IsTxRecordProcessed(ctx, hash) {
		return true, hash
	}

	return false, hash
}
