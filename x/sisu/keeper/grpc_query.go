package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type GrpcQuerier struct {
	keeper Keeper
}

func NewGrpcQuerier(keeper Keeper) *GrpcQuerier {
	return &GrpcQuerier{keeper: keeper}
}

func (k *GrpcQuerier) AllPubKeys(goCtx context.Context, req *types.QueryAllPubKeysRequest) (*types.QueryAllPubKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	log.Verbose("Fetching all pub keys.")

	allPubKeys := k.keeper.GetAllKeygenPubkeys(ctx)

	return &types.QueryAllPubKeysResponse{
		Pubkeys: allPubKeys,
	}, nil
}

func (k *GrpcQuerier) QueryContract(goCtx context.Context, req *types.QueryContractRequest) (*types.QueryContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	contract := k.keeper.GetContract(ctx, req.Chain, req.Hash, false)
	if contract == nil {
		return nil, fmt.Errorf("cannot find contract on chain %s and hash %s", req.Chain, req.Hash)
	}

	return &types.QueryContractResponse{
		Contract: contract,
	}, nil
}

func (k *GrpcQuerier) QueryToken(goCtx context.Context, req *types.QueryTokenRequest) (*types.QueryTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	m := k.keeper.GetTokens(ctx, []string{req.Id})
	if m[req.Id] == nil {
		return nil, fmt.Errorf("cannot find token %s", req.Id)
	}

	return &types.QueryTokenResponse{
		Token: m[req.Id],
	}, nil
}
