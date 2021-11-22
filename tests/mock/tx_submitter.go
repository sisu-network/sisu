package mock

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
)

// Make sure struct implement interface at compile-time
var _ common.TxSubmit = (*TxSubmitter)(nil)

type TxSubmitter struct {
	SubmitEThTxFunc   func(data []byte) error
	SubmitMessageFunc func(msg sdk.Msg) error
}

func (t TxSubmitter) SubmitEThTx(data []byte) error {
	if t.SubmitEThTxFunc == nil {
		return nil
	}

	return t.SubmitEThTxFunc(data)
}

func (t TxSubmitter) SubmitMessage(msg sdk.Msg) error {
	if t.SubmitMessageFunc == nil {
		return nil
	}

	return t.SubmitMessageFunc(msg)
}
