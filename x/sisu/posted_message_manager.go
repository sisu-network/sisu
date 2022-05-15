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
}

type DefaultPostedMessageManager struct {
	keeper     keeper.Keeper
	valManager ValidatorManager
}

func NewPostedMessageManager(keeper keeper.Keeper, valManager ValidatorManager) *DefaultPostedMessageManager {
	return &DefaultPostedMessageManager{
		keeper:     keeper,
		valManager: valManager,
	}
}

func (m *DefaultPostedMessageManager) ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte) {
	hash, signer, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Error("failed to get tx hash, err = ", err)
		return false, hash
	}

	m.keeper.SaveTxRecord(ctx, hash, signer)
	if m.valManager.HasConsensus(ctx, hash) && !m.keeper.IsTxRecordProcessed(ctx, hash) {
		return true, hash
	}

	return false, hash
}
