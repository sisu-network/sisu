package keeper

import (
	"context"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (k *Keeper) AllPubKeys(ctx context.Context, req *types.QueryAllPubKeysRequest) (*types.QueryAllPubKeysResponse, error) {
	utils.LogVerbose("Fetching all pub keys.")

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	pubKeys := k.GetAllPubKeys(sdkCtx)

	return &types.QueryAllPubKeysResponse{
		Pubkeys: pubKeys,
	}, nil
}
