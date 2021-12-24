package keeper

import (
	"fmt"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (k *DefaultKeeper) getContractKey(chain string, hash string) []byte {
	// chain + hash
	return []byte(fmt.Sprintf("%s__%s", chain, hash))
}

func (k *DefaultKeeper) getContractByteCodeKey(chain string, hash string) []byte {
	// chain + hash
	return []byte(fmt.Sprintf("%s__%s", chain, hash))
}

func (k *DefaultKeeper) SaveContracts(ctx sdk.Context, msgs []*types.Contract) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	byteCodeStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)

	for _, msg := range msgs {
		bz, err := msg.Marshal()
		if err != nil {
			log.Error("Cannot marshal contract message, err = ", err)
			continue
		}

		// Save byte code into separate store since it's rarely read.
		copy := &types.Contract{}
		err = copy.Unmarshal(bz)
		if err != nil {
			log.Error("Cannot unmarshal contract message, err = ", err)
			continue
		}

		// Set bytecode to nil
		copy.ByteCodes = nil
		bz, err = copy.Marshal()
		if err != nil {
			log.Error("Cannot marshal contract copy, err = ", err)
			continue
		}

		contractKey := k.getContractKey(msg.Chain, msg.Hash)
		byteCodeKey := k.getContractByteCodeKey(msg.Chain, msg.Hash)

		contractStore.Set(contractKey, bz)
		byteCodeStore.Set(byteCodeKey, msg.ByteCodes)
	}
}

func (k *DefaultKeeper) GetPendingContracts(ctx sdk.Context, chain string) []*types.Contract {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	byteCodeStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)

	contracts := make([]*types.Contract, 0)

	iter := contractStore.Iterator([]byte(fmt.Sprintf("%s__", chain)), nil)
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		bz := iter.Value()

		contract := &types.Contract{}
		err := contract.Unmarshal(bz)
		if err != nil {
			log.Error("Cannot unmarshal contract bytes")
			continue
		}

		if contract.Status != "" {
			continue
		}

		bz = byteCodeStore.Get(key)
		contract.ByteCodes = bz

		contracts = append(contracts, contract)
	}

	return contracts
}
