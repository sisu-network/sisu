package keeper

import (
	"context"
	"fmt"

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

func (k *DefaultKeeper) QueryContract(ctx context.Context, req *types.QueryContractRequest) (*types.QueryContractResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), prefixContract)

	contract := getContract(store, nil, req.Chain, req.Hash)
	if contract == nil {
		return nil, fmt.Errorf("cannot find contract on chain %s and hash %s", req.Chain, req.Hash)
	}

	return &types.QueryContractResponse{Contract: contract}, nil
}
