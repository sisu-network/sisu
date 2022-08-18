package cardano

import (
	"fmt"

	"github.com/echovl/cardano-go"
)

type NotEnoughBalanceErr struct {
	total, balance *cardano.Value
}

func NewNotEnoughBalanceErr(total, balance *cardano.Value) *NotEnoughBalanceErr {
	return &NotEnoughBalanceErr{
		total:   total,
		balance: balance,
	}
}

func (e *NotEnoughBalanceErr) Error() string {
	return fmt.Sprintf("Not enough balance, %v > %v", e.total, e.balance)
}
