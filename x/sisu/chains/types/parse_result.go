package types

import "github.com/sisu-network/sisu/x/sisu/types"

type ParseResult struct {
	Method       MethodType
	TransferOuts []*types.TransferDetails
	Error        error
}
