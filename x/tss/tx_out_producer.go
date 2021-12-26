package tss

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/db"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
	tsstypes "github.com/sisu-network/sisu/x/tss/types"
)

// This structs produces transaction output based on input. For a given tx input, this struct
// produces a list (could contain only one element) of transaction output.
type TxOutputProducer interface {
	AddKeyAddress(ctx sdk.Context, chain, addr string) error
	GetTxOuts(ctx sdk.Context, height int64, tx *types.ObservedTx) ([]*tssTypes.TxOut, []*tssTypes.TxOutEntity)
	SaveContractsToDeploy(chain string)
}

type DefaultTxOutputProducer struct {
	// List of key addresses in all eth based chain.
	// Map from: chain -> address -> bool.
	ethKeyAddrs map[string]map[string]bool

	worldState WorldState
	keeper     keeper.Keeper
	appKeys    common.AppKeys
	db         db.Database
	tssConfig  config.TssConfig
}

func NewTxOutputProducer(worldState WorldState, keeper keeper.Keeper, appKeys common.AppKeys, db db.Database, tssConfig config.TssConfig) TxOutputProducer {
	return &DefaultTxOutputProducer{
		worldState: worldState,
		keeper:     keeper,
		appKeys:    appKeys,
		tssConfig:  tssConfig,
		db:         db,
	}
}

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, height int64, tx *types.ObservedTx) ([]*tssTypes.TxOut, []*tssTypes.TxOutEntity) {
	outMsgs := make([]*tssTypes.TxOut, 0)
	outEntities := make([]*tssTypes.TxOutEntity, 0)
	var err error

	if libchain.IsETHBasedChain(tx.Chain) {
		log.Info("Getting tx out for chain ", tx.Chain)
		outMsgs, outEntities, err = p.getEthResponse(ctx, height, tx)

		if err != nil {
			log.Error("Cannot get response for an eth tx, err = ", err)
		}
	}

	return outMsgs, outEntities
}

func (p *DefaultTxOutputProducer) getEthKeyAddrs(ctx sdk.Context) (map[string]map[string]bool, error) {
	if p.ethKeyAddrs != nil {
		return p.ethKeyAddrs, nil
	}

	var err error
	p.ethKeyAddrs, err = p.keeper.GetAllEthKeyAddrs(ctx)
	if err != nil {
		return nil, err
	}

	return p.ethKeyAddrs, nil
}

func (p *DefaultTxOutputProducer) AddKeyAddress(ctx sdk.Context, chain, addr string) error {
	keyAddrs, err := p.getEthKeyAddrs(ctx)
	if err != nil {
		return err
	}

	m := keyAddrs[chain]
	if m == nil {
		m = make(map[string]bool)
	}
	m[addr] = true
	p.ethKeyAddrs[chain] = m

	return p.keeper.SaveEthKeyAddrs(ctx, chain, m)
}

// Get ETH out from an observed tx. Only do this if this is a validator node.
func (p *DefaultTxOutputProducer) getEthResponse(ctx sdk.Context, height int64, tx *types.ObservedTx) ([]*tsstypes.TxOut, []*tssTypes.TxOutEntity, error) {
	ethTx := &ethTypes.Transaction{}

	err := ethTx.UnmarshalBinary(tx.Serialized)
	if err != nil {
		log.Error("Failed to unmarshall eth tx. err =", err)
		return nil, nil, err
	}

	outMsgs := make([]*tssTypes.TxOut, 0)
	outEntities := make([]*tssTypes.TxOutEntity, 0)

	// 1. Check if this is a transaction sent to our key address. If this is true, it's likely a tx
	// that funds our account.
	if ethTx.To() != nil && p.db.IsChainKeyAddress(libchain.KEY_TYPE_ECDSA, ethTx.To().String()) {
		contracts := p.db.GetPendingDeployContracts(tx.Chain)
		log.Verbose("len(contracts) = ", len(contracts))

		if len(contracts) > 0 {
			// TODO: Check balance required to deploy all these contracts.

			// Get the list of deploy transactions. Those txs need to posted and verified (by validators)
			// to the Sisu chain.
			outEthTxs := p.checkEthDeployContract(ctx, height, tx.Chain, contracts)

			for i, outTx := range outEthTxs {
				bz, err := outTx.MarshalBinary()
				if err != nil {
					log.Error("Cannot marshall binary", err)
					continue
				}

				outMsg := tsstypes.NewMsgTxOut(
					p.appKeys.GetSignerAddress().String(),
					tx.BlockHeight,
					tx.Chain,
					tx.TxHash,
					tx.Chain,
					bz,
				)

				outMsgs = append(outMsgs, outMsg)

				outEntity := tssTypes.TxOutToEntity(outMsg)
				outEntity.ContractHash = contracts[i].Hash
				outEntities = append(outEntities, outEntity)
			}

			if len(outMsgs) > 0 {
				return outMsgs, outEntities, nil
			}
		}
	}

	// 2. Check if this is a tx sent to one of our contracts.
	if ethTx.To() != nil && len(ethTx.Data()) >= 4 {
		// TODO: compare method name to trigger corresponding function
		responseTx, err := p.processERC20TransferIn(ethTx)
		if err == nil {
			outMsg := tsstypes.NewMsgTxOut(
				p.appKeys.GetSignerAddress().String(),
				tx.BlockHeight,
				tx.Chain,
				tx.TxHash,
				responseTx.OutChain, // Could be different chain
				responseTx.RawBytes,
			)

			outMsgs = append(outMsgs, outMsg)

			outEntity := tssTypes.TxOutToEntity(outMsg)
			outEntities = append(outEntities, outEntity)
		}
	}

	// 3. Check other types of transaction.

	return outMsgs, outEntities, nil
}

// Check if we can deploy contract after seeing some ETH being sent to our ethereum address.
func (p *DefaultTxOutputProducer) checkEthDeployContract(ctx sdk.Context, height int64, chain string, contracts []*tsstypes.ContractEntity) []*ethTypes.Transaction {
	txs := make([]*ethTypes.Transaction, 0)

	for _, contract := range contracts {
		nonce := p.worldState.UseAndIncreaseNonce(chain)
		log.Verbose("nonce for deploying contract:", nonce)
		if nonce < 0 {
			log.Error("cannot get nonce for contract")
			continue
		}
		rawTx := p.getContractTx(contract, nonce)
		if rawTx == nil {
			log.Warn("raw Tx is nil")
			continue
		}

		txs = append(txs, rawTx)
	}

	// Update all contracts to "deploying" state.
	p.db.UpdateContractsStatus(contracts, tsstypes.ContractStateDeploying)

	return txs
}

// Save a list of contracts that are pending to be deployed. This is often call after a key is
// generated for a chain. We cannot deploy immediately after key generation because we don't have
// enough balance in the account.
func (p *DefaultTxOutputProducer) SaveContractsToDeploy(chain string) {
	if libchain.IsETHBasedChain(chain) {
		contracts := make([]*types.ContractEntity, 0, len(SupportedContracts))

		for name, c := range SupportedContracts {
			contract := &types.ContractEntity{
				Chain:    chain,
				Hash:     c.AbiHash,
				ByteCode: ecommon.FromHex(c.Bin),
				Name:     name,
			}

			contracts = append(contracts, contract)
		}

		p.db.InsertContracts(contracts)
	}
}

func (p *DefaultTxOutputProducer) getContractTx(contract *tsstypes.ContractEntity, nonce int64) *ethTypes.Transaction {
	erc20gw := SupportedContracts[ContractErc20Gateway]
	switch contract.Hash {
	case erc20gw.AbiHash:
		// This is erc20gw contract.
		parsedAbi, err := abi.JSON(strings.NewReader(erc20gw.AbiString))
		if err != nil {
			log.Error("cannot parse erc20 abi. abi = ", erc20gw.AbiString, "err =", err)
			return nil
		}

		// Get all allowed chains
		supportedChains := make([]string, 0)
		for chain := range p.tssConfig.SupportedChains {
			if chain != contract.Chain {
				supportedChains = append(supportedChains, chain)
			}
		}

		log.Info("Allowed chains for chain", contract.Chain, "are: ", supportedChains)

		input, err := parsedAbi.Pack("", supportedChains)
		if err != nil {
			log.Error("cannot pack supportedChains, err =", err)
			return nil
		}

		byteCode := ecommon.FromHex(erc20gw.Bin)
		input = append(byteCode, input...)

		rawTx := ethTypes.NewTx(&ethTypes.AccessListTx{
			Nonce:    uint64(nonce),
			GasPrice: p.getGasPrice(contract.Chain),
			Gas:      p.getGasLimit(contract.Chain),
			Value:    big.NewInt(0),
			Data:     input,
		})
		return rawTx
	}

	return nil
}

func (p *DefaultTxOutputProducer) getGasLimit(chain string) uint64 {
	// TODO: Make this dependent on different chains.
	return uint64(8_000_000)
}

func (p *DefaultTxOutputProducer) getGasPrice(chain string) *big.Int {
	// TODO: Make this dependent on different chains.
	//switch chain {
	//case "eth-ropsten":
	//	return big.NewInt(1700000000)
	//case "eth-binance-testnet":
	//	return big.NewInt(10000000000) // 10 Gwei
	//}
	return big.NewInt(400000000000) // 10 Gwei
}
