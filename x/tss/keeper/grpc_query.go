package keeper

import (
	"context"

	sdk "github.com/sisu-network/cosmos-sdk/types"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (k *DefaultKeeper) AllPubKeys(ctx context.Context, req *types.QueryAllPubKeysRequest) (*types.QueryAllPubKeysResponse, error) {
	log.Verbose("Fetching all pub keys.")

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), prefixKeygen)

	return &types.QueryAllPubKeysResponse{
		Pubkeys: getAllKeygenPubkeys(store),
	}, nil
}
