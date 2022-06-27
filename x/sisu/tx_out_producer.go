package sisu

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	etypes "github.com/sisu-network/deyes/types"
	hutils "github.com/sisu-network/dheart/utils"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	scardano "github.com/sisu-network/sisu/x/sisu/cardano"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
)

// This structs produces transaction output based on input. For a given tx input, this struct
// produces a list (could contain only one element) of transaction output.
type TxOutputProducer interface {
	GetTxOuts(ctx sdk.Context, height int64, tx *types.TxIn) []*types.TxOutWithSigner

	PauseContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)

	ResumeContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)

	ContractChangeOwnership(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutWithSigner, error)

	ContractSetLiquidPoolAddress(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutWithSigner, error)

	ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutWithSigner, error)
}

type DefaultTxOutputProducer struct {
	// List of key addresses in all eth based chain.
	// Map from: chain -> address -> bool.
	ethKeyAddrs map[string]map[string]bool

	worldState     world.WorldState
	appKeys        common.AppKeys
	keeper         keeper.Keeper
	tssConfig      config.TssConfig
	cardanoConfig  config.CardanoConfig
	cardanoNetwork cardano.Network

	// Only use for cardano chain
	cardanoNode cardano.Node
}

func NewTxOutputProducer(worldState world.WorldState, appKeys common.AppKeys,
	keeper keeper.Keeper, tssConfig config.TssConfig, cardanoConfig config.CardanoConfig, cardanoNode cardano.Node) TxOutputProducer {
	return &DefaultTxOutputProducer{
		keeper:         keeper,
		worldState:     worldState,
		appKeys:        appKeys,
		tssConfig:      tssConfig,
		cardanoNetwork: cardanoConfig.GetCardanoNetwork(),
		cardanoNode:    cardanoNode,
	}
}

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, height int64, tx *types.TxIn) []*types.TxOutWithSigner {
	outMsgs := make([]*types.TxOutWithSigner, 0)
	var err error

	if libchain.IsETHBasedChain(tx.Chain) {
		log.Info("Getting tx out for chain ", tx.Chain)
		outMsgs, err = p.getEthResponse(ctx, height, tx)

		if err != nil {
			log.Error("Cannot get response for an eth tx, err = ", err)
			return nil
		}
	}

	if libchain.IsCardanoChain(tx.Chain) {
		log.Info("Found tx in request from Cardano to Ethereum")
		outMsg, err := p.extractCardanoTxIn(ctx, tx)
		if err != nil {
			return nil
		}

		outMsgs = append(outMsgs, outMsg)

	}

	return outMsgs
}

// Get ETH out from an observed tx. Only do this if this is a validator node.
func (p *DefaultTxOutputProducer) getEthResponse(ctx sdk.Context, height int64, tx *types.TxIn) ([]*types.TxOutWithSigner, error) {
	ethTx := &ethTypes.Transaction{}

	err := ethTx.UnmarshalBinary(tx.Serialized)
	if err != nil {
		log.Error("Failed to unmarshall eth tx. err =", err)
		return nil, err
	}

	outMsgs := make([]*types.TxOutWithSigner, 0)

	log.Verbose("ethTx.To() = ", ethTx.To())

	// 1. Check if this is a transaction sent to our key address. If this is true, it's likely a tx
	// that funds our account.
	if ethTx.To() != nil && p.keeper.IsKeygenAddress(ctx, libchain.KEY_TYPE_ECDSA, ethTx.To().String()) {
		contracts := p.keeper.GetPendingContracts(ctx, tx.Chain)
		log.Verbose("len(contracts) = ", len(contracts))

		if len(contracts) > 0 {
			// TODO: Check balance required to deploy all these contracts.

			// Get the list of deploy transactions. Those txs need to posted and verified (by validators)
			// to the Sisu chain.
			outEthTxs := p.getEthContractDeploymentTx(ctx, height, tx.Chain, contracts)

			for i, outTx := range outEthTxs {
				bz, err := outTx.MarshalBinary()
				if err != nil {
					log.Error("Cannot marshall binary", err)
					continue
				}

				outMsg := types.NewMsgTxOutWithSigner(
					p.appKeys.GetSignerAddress().String(),
					types.TxOutType_CONTRACT_DEPLOYMENT,
					tx.BlockHeight,
					tx.Chain,
					tx.TxHash,
					tx.Chain,
					outTx.Hash().String(),
					bz,
					contracts[i].Hash,
				)

				log.Verbose("ETH Tx Out hash = ", outTx.Hash().String(), " on chain ", tx.Chain)

				outMsgs = append(outMsgs, outMsg)
			}

			if len(outMsgs) > 0 {
				return outMsgs, nil
			}
		}
	}

	log.Verbose("len(ethTx.Data()) ", len(ethTx.Data()))

	// 2. Check if this is a tx sent to one of our contracts.
	if ethTx.To() != nil &&
		p.keeper.IsContractExistedAtAddress(ctx, tx.Chain, ethTx.To().String()) && len(ethTx.Data()) >= 4 {
		data, err := parseEthTransferOut(ethTx, tx.Chain, p.worldState)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		if libchain.IsETHBasedChain(data.destChain) {
			// This is a swap from ETH -> ETH
			responseTx, err := p.buildERC20TransferIn(ctx, data.token, data.tokenAddr, ecommon.HexToAddress(data.recipient),
				data.amount, data.destChain)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			outMsg := types.NewMsgTxOutWithSigner(
				p.appKeys.GetSignerAddress().String(),
				types.TxOutType_TRANSFER_OUT,
				tx.BlockHeight,
				tx.Chain,
				tx.TxHash,
				responseTx.OutChain, // Could be different chain
				responseTx.EthTx.Hash().String(),
				responseTx.RawBytes,
				"",
			)

			outMsgs = append(outMsgs, outMsg)
		}

		if libchain.IsCardanoChain(data.destChain) {
			// This is a swap from ETH -> Cardano
			// Convert the ETH big.Int amount to lovelace. Most ERC20 has 18 decimals while lovelace has
			// only 6 decimals.
			multiAssetAmt, err := utils.SourceAmountToLovelace(tx.Chain, data.amount)
			if err != nil {
				return nil, err
			}
			log.Verbosef("data.amount = %v, multiAssetAmt = %v", data.amount, multiAssetAmt)

			// In real, this transaction transfers at least <1 ADA + additional tx fee>
			cardanoTx, err := p.getCardanoTx(ctx, data, multiAssetAmt.Uint64())

			if err != nil {
				return nil, err
			}

			cardanoTxHash, err := cardanoTx.Hash()
			if err != nil {
				return nil, err
			}

			bz, err := cardanoTx.MarshalCBOR()
			if err != nil {
				return nil, err
			}

			outMsg := types.NewMsgTxOutWithSigner(
				p.appKeys.GetSignerAddress().String(),
				types.TxOutType_TRANSFER_OUT,
				tx.BlockHeight,
				tx.Chain,
				tx.TxHash,
				data.destChain,
				cardanoTxHash.String(),
				bz,
				"",
			)

			outMsgs = append(outMsgs, outMsg)
		}
	}

	// 3. Check other types of transaction.

	return outMsgs, nil
}

// In Cardano chain, transferring multi-asset required at least 1 ADA (10^6 lovelace)
func (p *DefaultTxOutputProducer) getCardanoTx(ctx sdk.Context, data *transferOutData, assetAmount uint64) (*cardano.Tx, error) {
	pubkey := p.keeper.GetKeygenPubkey(ctx, libchain.KEY_TYPE_EDDSA)
	senderAddr := hutils.GetAddressFromCardanoPubkey(pubkey)
	log.Debug("cardano sender address = ", senderAddr.String())

	receiverAddr, err := cardano.NewAddress(data.recipient)
	if err != nil {
		log.Error("error when parsing receiver addr: ", err)
		return nil, err
	}

	multiAsset, err := scardano.GetCardanoMultiAsset(data.destChain, data.token, assetAmount)
	if err != nil {
		return nil, err
	}

	adaPrice, err := p.worldState.GetNativeTokenPriceForChain(data.destChain)
	if err != nil {
		log.Error("error when getting ada price: ", err)
		return nil, err
	}
	log.Debug("ADA price = ", adaPrice)

	// We need at least 1 ada to send multi assets.
	tx, err := scardano.BuildTx(p.cardanoNode, p.cardanoNetwork, senderAddr, receiverAddr,
		cardano.NewValueWithAssets(cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()), multiAsset), nil,
		adaPrice, data.token, data.destChain, assetAmount)

	if err != nil {
		log.Error("error when building tx: ", err)
		return nil, err
	}

	for _, i := range tx.Body.Inputs {
		log.Debugf("tx input = %v\n", i)
	}

	for _, o := range tx.Body.Outputs {
		log.Debugf("tx output = %v\n", o)
	}

	log.Debug("tx fee = ", tx.Body.Fee)

	return tx, nil
}

// Check if we can deploy contract after seeing some ETH being sent to our ethereum address.
func (p *DefaultTxOutputProducer) getEthContractDeploymentTx(ctx sdk.Context, height int64, chain string, contracts []*types.Contract) []*ethTypes.Transaction {
	txs := make([]*ethTypes.Transaction, 0)

	for _, contract := range contracts {
		nonce := p.worldState.UseAndIncreaseNonce(ctx, chain)
		log.Verbose("nonce for deploying contract:", nonce, " on chain ", chain)
		if nonce < 0 {
			log.Error("cannot get nonce for contract")
			continue
		}

		rawTx := p.getContractTx(ctx, contract, nonce)
		if rawTx == nil {
			log.Warn("raw Tx is nil")
			continue
		}

		txs = append(txs, rawTx)
	}

	return txs
}

func (p *DefaultTxOutputProducer) getContractTx(ctx sdk.Context, contract *types.Contract, nonce int64) *ethTypes.Transaction {
	erc20 := SupportedContracts[ContractErc20Gateway]
	switch contract.Hash {
	case erc20.AbiHash:
		// This is erc20gw contract.
		parsedAbi, err := abi.JSON(strings.NewReader(erc20.AbiString))
		if err != nil {
			log.Error("cannot parse erc20 abi. abi = ", erc20.AbiString, "err =", err)
			return nil
		}

		// Get all allowed chains
		supportedChains := make([]string, 0)
		for chain := range p.tssConfig.SupportedChains {
			if chain != contract.Chain {
				supportedChains = append(supportedChains, chain)
			}
		}

		sort.Strings(supportedChains)

		log.Info("Allowed chains for chain ", contract.Chain, " are: ", supportedChains)

		lp := p.keeper.GetLiquidity(ctx, contract.Chain)
		if lp == nil {
			log.Warn("Lp is nil for chain ", contract.Chain)
			return nil
		}

		log.Infof("Liquidity pool addr for chain %s is %s", contract.Chain, lp.Address)
		input, err := parsedAbi.Pack("", supportedChains, ecommon.HexToAddress(lp.Address))
		if err != nil {
			log.Error("cannot pack supportedChains, err =", err)
			return nil
		}

		byteCode := ecommon.FromHex(erc20.Bin)
		input = append(byteCode, input...)
		chain := p.keeper.GetChain(ctx, contract.Chain)
		if chain == nil {
			log.Error("getContractTx: chain is nil with id ", contract.Chain)
			return nil
		}

		gasPrice := chain.GasPrice
		if gasPrice <= 0 {
			gasPrice = p.getDefaultGasPrice(contract.Chain).Int64()
		}
		gasLimit := p.getGasLimit(contract.Chain)

		log.Verbose("Gas price = ", gasPrice, " on chain ", contract.Chain)
		log.Verbose("gasLimit = ", gasLimit, " on chain ", contract.Chain)

		rawTx := ethTypes.NewContractCreation(
			uint64(nonce),
			big.NewInt(0),
			gasLimit,
			big.NewInt(gasPrice),
			input,
		)

		return rawTx
	}

	return nil
}

func (p *DefaultTxOutputProducer) extractCardanoTxIn(ctx sdk.Context, tx *types.TxIn) (*types.TxOutWithSigner, error) {
	txIn := &etypes.CardanoTxInItem{}
	if err := json.Unmarshal(tx.Serialized, txIn); err != nil {
		log.Error("error when marshaling cardano tx in item: ", err)
		return nil, err
	}

	extraInfo := txIn.Metadata
	tokenName := txIn.Asset
	if tokenName != "ADA" {
		tokenName = tokenName[5:] // Remove the WRAP_ prefix
	}
	if len(tokenName) == 0 {
		return nil, fmt.Errorf("Invalid token: %s", tokenName)
	}

	tokens := p.keeper.GetTokens(ctx, []string{tokenName})
	token := tokens[tokenName]
	if token == nil {
		return nil, fmt.Errorf("Cannot find token in the keeper")
	}
	dstTokenAddress := token.GetAddressForChain(txIn.Metadata.Chain)

	amount := new(big.Int).SetUint64(txIn.Amount)
	recipient := ecommon.HexToAddress(extraInfo.Recipient)
	tokenAddr := ecommon.HexToAddress(dstTokenAddress)

	responseTx, err := p.buildERC20TransferIn(ctx, token, tokenAddr, recipient,
		utils.LovelaceToETHTokens(amount), extraInfo.Chain)
	if err != nil {
		return nil, err
	}

	outMsg := types.NewMsgTxOutWithSigner(
		p.appKeys.GetSignerAddress().String(),
		types.TxOutType_TRANSFER_OUT,
		tx.BlockHeight,
		tx.Chain,
		tx.TxHash,
		responseTx.OutChain,
		responseTx.EthTx.Hash().String(),
		responseTx.RawBytes,
		"",
	)

	return outMsg, nil
}

func (p *DefaultTxOutputProducer) getGasLimit(chain string) uint64 {
	// TODO: Make this dependent on different chains.
	return uint64(8_000_000)
}

// @Deprecated
func (p *DefaultTxOutputProducer) getDefaultGasPrice(chain string) *big.Int {
	// TODO: Make this dependent on different chains.
	switch chain {
	case "ganache1":
		return big.NewInt(2_000_000_000)
	case "ganache2":
		return big.NewInt(2_000_000_000)
	case "ropsten-testnet":
		return big.NewInt(4_000_000_000)
	case "binance-testnet":
		return big.NewInt(18_000_000_000)
	case "polygon-testnet":
		return big.NewInt(7_000_000_000)
	case "xdai":
		return big.NewInt(2_000_000_000)
	case "goerli-testnet":
		return big.NewInt(1_500_000_000)
	case "eth":
		return big.NewInt(70_000_000_000)
	case "arbitrum-testnet":
		return big.NewInt(50_000_000)
	case "fantom-testnet":
		return big.NewInt(75_000_000_000)
	}
	return big.NewInt(100_000_000_000)
}
