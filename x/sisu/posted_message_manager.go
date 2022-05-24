package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

var _ PostedMessageManager = (*DefaultPostedMessageManager)(nil)

//go:generate mockgen -source=./x/sisu/posted_message_manager.go -destination=./tests/mock/x/sisu/posted_message_manager.go -package=mock
type PostedMessageManager interface {
	ProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte, error)
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

// ProcessMsg returns true if this message should be processed
// If this message is processed before, do nothing
// else if this message reached consensus, payback the slash tokens to voters
func (m *DefaultPostedMessageManager) ProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte, error) {
	hash, signer, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Error("failed to get tx hash, err = ", err)
		return false, nil, err
	}

	if m.keeper.IsTxRecordProcessed(ctx, hash) {
		return false, hash, nil
	}

	bepAddr, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return false, hash, err
	}

	if err := m.keeper.IncSlashToken(ctx, getSlashToken(msg), bepAddr); err != nil {
		return false, nil, err
	}

	voters := m.keeper.GetVotersInAccAddress(ctx, hash)
	for _, v := range voters {
		b, err := m.keeper.GetSlashToken(ctx, v)
		if err != nil {
			log.Debug("error when getting slash token")
			continue
		}

		log.Debugf("AFTER INC: addr %s with slash point = %s", v.String(), b)
	}

	m.keeper.SaveTxRecord(ctx, hash, signer)

	if m.valManager.HasConsensus(ctx, hash) {
		// payback slash token to voters
		voterAddrs := m.keeper.GetVotersInAccAddress(ctx, hash)
		if err := m.keeper.DecSlashToken(ctx, getSlashToken(msg), voterAddrs...); err != nil {
			return false, hash, err
		}

		for _, v := range voterAddrs {
			b, err := m.keeper.GetSlashToken(ctx, v)
			if err != nil {
				log.Debug("error when getting slash token")
				continue
			}

			log.Debugf("AFTER DEC: addr %s with slash point = %d", v.String(), b)
		}
		return true, hash, nil
	}

	return false, hash, nil
}

func getSlashToken(msg sdk.Msg) int64 {
	return int64(len(msg.String()))
}
