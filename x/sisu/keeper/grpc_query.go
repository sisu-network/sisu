package keeper

import (
	"context"
	"fmt"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type GrpcQuerier struct {
	storage Storage
}

func NewGrpcQuerier(storage Storage) *GrpcQuerier {
	return &GrpcQuerier{storage: storage}
}

func (k *GrpcQuerier) AllPubKeys(ctx context.Context, req *types.QueryAllPubKeysRequest) (*types.QueryAllPubKeysResponse, error) {
	log.Verbose("Fetching all pub keys.")

	allPubKeys := k.storage.GetAllKeygenPubkeys()

	return &types.QueryAllPubKeysResponse{
		Pubkeys: allPubKeys,
	}, nil
}

func (k *GrpcQuerier) QueryContract(ctx context.Context, req *types.QueryContractRequest) (*types.QueryContractResponse, error) {
	contract := k.storage.GetContract(req.Chain, req.Hash, false)
	if contract == nil {
		return nil, fmt.Errorf("cannot find contract on chain %s and hash %s", req.Chain, req.Hash)
	}

	return &types.QueryContractResponse{
		Contract: contract,
	}, nil
}

func (k *GrpcQuerier) QueryToken(ctx context.Context, req *types.QueryTokenRequest) (*types.QueryTokenResponse, error) {
	m := k.storage.GetTokens([]string{req.Id})
	if m[req.Id] == nil {
		return nil, fmt.Errorf("cannot find token %s", req.Id)
	}

	return &types.QueryTokenResponse{
		Token: m[req.Id],
	}, nil
}
