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

func (k *GrpcQuerier) AllPubKeys(goCtx context.Context, req *types.QueryAllPubKeysRequest) (
	*types.QueryAllPubKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	log.Verbose("Fetching all pub keys.")

	allPubKeys := k.keeper.GetAllKeygenPubkeys(ctx)

	return &types.QueryAllPubKeysResponse{
		Pubkeys: allPubKeys,
	}, nil
}

func (k *GrpcQuerier) QueryVault(goCtx context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	v := k.keeper.GetVault(ctx, req.Chain, req.Token)
	if v == nil {
		return nil, fmt.Errorf("cannot find contract on chain %s", req.Chain)
	}

	return &types.QueryVaultResponse{
		Vault: v,
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

func (k *GrpcQuerier) QueryChain(goCtx context.Context, req *types.QueryChainRequest) (*types.QueryChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.keeper.GetChain(ctx, req.Chain)
	if chain == nil {
		return nil, fmt.Errorf("cannot find chain %s", req.Chain)
	}

	return &types.QueryChainResponse{
		Chain: chain,
	}, nil
}
