package sisu

import "fmt"

/////

var (
	ErrInsufficientFund = fmt.Errorf("insufficient fund")
	ErrInvalidMessage   = fmt.Errorf("invalid message")
)

///// Chain not found

type ErrChainNotFound struct {
	chain string
}

func NewErrChainNotFound(chain string) ErrChainNotFound {
	return ErrChainNotFound{
		chain: chain,
	}
}

func (e ErrChainNotFound) Error() string {
	return fmt.Sprintf("chain %s not found", e.chain)
}

///// Token not found

type ErrTokenNotFound struct {
	token string
}

func NewErrTokenNotFound(token string) ErrTokenNotFound {
	return ErrTokenNotFound{
		token: token,
	}
}

func (e ErrTokenNotFound) Error() string {
	return fmt.Sprintf("token %s not found", e.token)
}
