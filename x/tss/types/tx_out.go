package types

import (
	"github.com/sisu-network/sisu/utils"
)

func (msg *TxOut) GetHash() string {
	return utils.KeccakHash32(msg.OutChain + string(msg.OutBytes))
}
