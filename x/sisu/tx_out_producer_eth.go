package sisu

import (
	"fmt"
	"math/big"

	ecommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *DefaultTxOutputProducer) buildERC20TransferIn(
	ctx sdk.Context,
	k keeper.Keeper,
	tokens []*types.Token,
	recipients []ethcommon.Address,
	amounts []*big.Int,
	destChain string,
) (*types.TxResponse, error) {
	targetContractName := ContractVault
	v := p.keeper.GetVault(ctx, destChain, "")
	if v == nil {
		return nil, fmt.Errorf("Cannot find vault for chain %s", destChain)
	}
	gw := v.Address
	if len(gw) == 0 {
		err := fmt.Errorf("cannot find gw address for type: %s on chain %s", targetContractName, destChain)
		log.Error(err)
		return nil, err
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	vaultInfo := SupportedContracts[targetContractName]

	chain := k.GetChain(ctx, destChain)
	if chain == nil {
		return nil, fmt.Errorf("Invalid chain: %s", chain)
	}

	gasPrice := big.NewInt(chain.GasPrice)
	if gasPrice.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("Gas price is non-positive: %s", gasPrice.String())
	}

	commissionRate := p.keeper.GetParams(ctx).CommissionRate
	if commissionRate < 0 || commissionRate > 10_000 {
		return nil, fmt.Errorf("Commission rate is invalid, rate = %d", commissionRate)
	}

	log.Debug("Gas price for swapping  = ", gasPrice)

	finalTokenAddrs := make([]ethcommon.Address, 0)
	finalRecipients := make([]ethcommon.Address, 0)
	finalAmounts := make([]*big.Int, 0)
	amountIns := make([]*big.Int, 0)
	gasPrices := make([]*big.Int, 0)

	for i := range amounts {
		amountOut := new(big.Int).Set(amounts[i])

		// Subtract commission rate
		amountOut = utils.SubtractCommissionRate(amountOut, commissionRate)

		price, ok := new(big.Int).SetString(tokens[i].Price, 10)
		if !ok {
			return nil, fmt.Errorf("invalid token price %s", tokens[i].Price)
		}
		if price.Cmp(utils.ZeroBigInt) == 0 {
			return nil, fmt.Errorf("token %s has price 0", tokens[i].Id)
		}

		gasPriceInToken, err := helper.GetChainGasCostInToken(ctx, k, tokens[i].Id, destChain, big.NewInt(80_000))
		if err != nil {
			log.Error("Cannot get gas cost in token, err = ", err)
			continue
		}

		if gasPriceInToken.Cmp(utils.ZeroBigInt) < 0 {
			log.Errorf("Gas price in token is negative: token id = %s", tokens[i].Id)
			gasPriceInToken = utils.ZeroBigInt
		}

		// Subtract gas price in token.
		amountOut.Sub(amountOut, gasPriceInToken)

		if amountOut.Cmp(utils.ZeroBigInt) < 0 {
			log.Error("Insufficient fund for transfer amountOut = ", amountOut, " gasPriceInToken = ", gasPriceInToken)
			continue
		}

		// Find the address of the token.
		var tokenAddr string
		for _, token := range tokens {
			for j, chain := range token.Chains {
				if chain == destChain {
					tokenAddr = token.Addresses[j]
					break
				}
			}
			if len(tokenAddr) > 0 {
				break
			}
		}
		if len(tokenAddr) == 0 {
			continue
		}

		finalTokenAddrs = append(finalTokenAddrs, ecommon.HexToAddress(tokenAddr))
		finalAmounts = append(finalAmounts, amountOut)
		finalRecipients = append(finalRecipients, recipients[i])
		amountIns = append(amountIns, amounts[i])
		gasPrices = append(gasPrices, gasPriceInToken)
	}

	if len(finalTokenAddrs) == 0 {
		return nil, fmt.Errorf("No txOut is produced (might be due to insufficient fund")
	}

	log.Verbosef("destChain: %s, gateway address on destChain: %s", destChain, gatewayAddress.String())
	for i := range finalTokenAddrs {
		log.Verbosef("tokenAddr: %s, recipient: %s, gasPriceInToken: %d, amountIn: %s, amountOut: %s",
			finalTokenAddrs[i], finalRecipients[i], gasPrices[i], amountIns[i].String(), finalAmounts[i].String(),
		)
	}

	var input []byte
	var err error
	if len(finalTokenAddrs) == 1 {
		input, err = vaultInfo.Abi.Pack(
			MethodTransferIn,
			finalTokenAddrs[0],
			finalRecipients[0],
			finalAmounts[0],
		)
	} else {
		input, err = vaultInfo.Abi.Pack(
			MethodTransferInMultiple,
			finalTokenAddrs,
			finalRecipients,
			finalAmounts,
		)
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	rawTx := ethTypes.NewTransaction(
		0,
		gatewayAddress,
		big.NewInt(0),
		uint64(100_000*len(recipients)), // 100k per swapping operation.
		gasPrice,
		input,
	)

	bz, err := rawTx.MarshalBinary()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &types.TxResponse{
		OutChain: destChain,
		EthTx:    rawTx,
		RawBytes: bz,
	}, nil
}
