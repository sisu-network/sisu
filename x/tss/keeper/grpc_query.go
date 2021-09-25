package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (k *Keeper) AllPubKeys(ctx context.Context, req *types.QueryAllPubKeysRequest) (*types.QueryAllPubKeysResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	pubKeys := k.GetAllPubKeys(sdkCtx)

	return &types.QueryAllPubKeysResponse{
		Pubkeys: pubKeys,
	}, nil
}
