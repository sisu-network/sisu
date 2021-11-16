package keeper

import (
	"context"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (k *DefaultKeeper) AllPubKeys(ctx context.Context, req *types.QueryAllPubKeysRequest) (*types.QueryAllPubKeysResponse, error) {
	log.Verbose("Fetching all pub keys.")

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	pubKeys := k.GetAllPubKeys(sdkCtx)

	return &types.QueryAllPubKeysResponse{
		Pubkeys: pubKeys,
	}, nil
}
