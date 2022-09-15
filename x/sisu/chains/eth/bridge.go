package eth

import (
	"fmt"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	ctypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type ethBridge struct {
	signer string
	chain  string
	keeper keeper.Keeper
}

func NewBridge(chain string, signer string, k keeper.Keeper) ctypes.Bridge {
	return &ethBridge{
		chain:  chain,
		signer: signer,
		keeper: k,
	}
}

func (b *ethBridge) ProcessTransfers(ctx sdk.Context, transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
	inHashes := make([]string, 0, len(transfers))
	tokens := make([]*types.Token, 0, len(transfers))
	recipients := make([]ethcommon.Address, 0, len(transfers))
	amounts := make([]*big.Int, 0, len(transfers))

	allTokens := b.keeper.GetAllTokens(ctx)
	for _, transfer := range transfers {
		token := allTokens[transfer.Token]
		if token == nil {
			log.Warn("cannot find token", transfer.Token)
			continue
		}

		amount, ok := new(big.Int).SetString(transfer.Amount, 10)
		if !ok {
			log.Warn("Cannot create big.Int value from amout ", transfer.Amount)
			continue
		}

		tokens = append(tokens, token)
		recipients = append(recipients, ethcommon.HexToAddress(transfer.ToRecipient))
		amounts = append(amounts, amount)
		inHashes = append(inHashes, transfer.Id)

		log.Verbosef("Processing transfer in: id = %s, recipient = %s, amount = %s, inHash = %s",
			token.Id, transfer.ToRecipient, amount, transfer.Id)
	}

	responseTx, err := b.buildERC20TransferIn(ctx, tokens, recipients, amounts)
	if err != nil {
		log.Error("Failed to build erc20 transfer in, err = ", err)
		return nil, err
	}

	outMsg := types.NewTxOutMsg(
		b.signer,
		types.TxOutType_TRANSFER_OUT,
		&types.TxOutContent{
			OutChain: b.chain,
			OutHash:  responseTx.EthTx.Hash().String(),
			OutBytes: responseTx.RawBytes,
		},
		&types.TxOutInput{
			TransferIds: inHashes,
		},
	)

	return []*types.TxOutMsg{outMsg}, nil
}

func (b *ethBridge) buildERC20TransferIn(
	ctx sdk.Context,
	tokens []*types.Token,
	recipients []ethcommon.Address,
	amounts []*big.Int,
) (*types.TxResponse, error) {
	targetContractName := ContractVault
	v := b.keeper.GetVault(ctx, b.chain)
	if v == nil {
		return nil, fmt.Errorf("Cannot find vault for chain %s", b.chain)
	}
	gw := v.Address
	if len(gw) == 0 {
		err := fmt.Errorf("cannot find gw address for type: %s on chain %s", targetContractName, b.chain)
		log.Error(err)
		return nil, err
	}

	gatewayAddress := ethcommon.HexToAddress(gw)
	vaultInfo := SupportedContracts[targetContractName]

	chain := b.keeper.GetChain(ctx, b.chain)
	if chain == nil {
		return nil, fmt.Errorf("Invalid chain: %s", chain)
	}

	gasPrice := big.NewInt(chain.GasPrice)
	if gasPrice.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("Gas price is non-positive: %s", gasPrice.String())
	}

	commissionRate := b.keeper.GetParams(ctx).CommissionRate
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

		gasPriceInToken, err := helper.GetChainGasCostInToken(ctx, b.keeper, tokens[i].Id, b.chain, big.NewInt(80_000))
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
				if chain == b.chain {
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

		finalTokenAddrs = append(finalTokenAddrs, ethcommon.HexToAddress(tokenAddr))
		finalAmounts = append(finalAmounts, amountOut)
		finalRecipients = append(finalRecipients, recipients[i])
		amountIns = append(amountIns, amounts[i])
		gasPrices = append(gasPrices, gasPriceInToken)
	}

	if len(finalTokenAddrs) == 0 {
		return nil, fmt.Errorf("No txOut is produced (might be due to insufficient fund")
	}

	log.Verbosef("destChain: %s, gateway address on destChain: %s", b.chain, gatewayAddress.String())
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

	rawTx := ethtypes.NewTransaction(
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
		OutChain: b.chain,
		EthTx:    rawTx,
		RawBytes: bz,
	}, nil
}
