package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CloneSdkContext clones a context. This function is often needed when we need a readOnly context
// and does not want to modify the original context.
func CloneSdkContext(ctx sdk.Context) sdk.Context {
	clone := sdk.Context{}
	cacheMS := ctx.MultiStore().CacheMultiStore()

	clone = sdk.NewContext(
		cacheMS, ctx.BlockHeader(), ctx.IsCheckTx(), nil,
	)

	return clone
}
