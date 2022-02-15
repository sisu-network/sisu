package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

type PostedMessageManager interface {
	ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte)
}

type DefaultPostedMessageManager struct {
	publicDb  keeper.Storage
	threshold int
}

func NewPostedMessageManager(publicDb keeper.Storage, threshold int) *DefaultPostedMessageManager {
	return &DefaultPostedMessageManager{
		publicDb:  publicDb,
		threshold: threshold,
	}
}

func (m *DefaultPostedMessageManager) ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte) {
	hash, signer, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Error("failed to get tx hash, err = ", err)
		return false, hash
	}

	count := m.publicDb.SaveTxRecord(hash, signer)
	if count >= m.threshold && !m.publicDb.IsTxRecordProcessed(hash) {
		return true, hash
	}

	return false, hash
}
