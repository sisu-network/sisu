package sisu

import (
	"math/big"

	scardano "github.com/sisu-network/sisu/x/sisu/chains/cardano"

	ethcommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"

	chainseth "github.com/sisu-network/sisu/x/sisu/chains/eth"
	chainstypes "github.com/sisu-network/sisu/x/sisu/chains/types"
)

// This structs produces transaction output based on input. For a given tx input, this struct
// produces a list (could contain only one element) of transaction output.
type TxOutputProducer interface {
	// GetTxOuts returns a list of TxOut message and a list of un-processed transfer out request that
	// needs to be processed next time.
	GetTxOuts(ctx sdk.Context, chain string, transfers []*types.Transfer) ([]*types.TxOutMsg, error)
}

type DefaultTxOutputProducer struct {
	signer    string
	keeper    keeper.Keeper
	txTracker TxTracker

	// Only use for cardano chain
	cardanoConfig  config.CardanoConfig
	cardanoNetwork cardano.Network
	cardanoClient  scardano.CardanoClient

	bridges map[string]chainstypes.Bridge
}

type transferInData struct {
	token     ethcommon.Address
	recipient string
	amount    *big.Int
}

func NewTxOutputProducer(appKeys common.AppKeys, keeper keeper.Keeper,
	cardanoConfig config.CardanoConfig,
	cardanoClient scardano.CardanoClient,
	txTracker TxTracker) TxOutputProducer {
	return &DefaultTxOutputProducer{
		signer:         appKeys.GetSignerAddress().String(),
		keeper:         keeper,
		txTracker:      txTracker,
		cardanoNetwork: cardanoConfig.GetCardanoNetwork(),
		cardanoClient:  cardanoClient,
		bridges:        make(map[string]chainstypes.Bridge),
	}
}

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, chain string,
	transfers []*types.Transfer) ([]*types.TxOutMsg, error) {

	if libchain.IsETHBasedChain(chain) {
		bridge := p.getBridge(chain)
		msgs, err := bridge.ProcessTransfers(ctx, transfers)
		if err != nil {
			return nil, err
		}

		return msgs, nil
	}

	if libchain.IsCardanoChain(chain) {
		msgs, err := p.processCardanoBatches(ctx, p.keeper, chain, transfers)
		if err != nil {
			return nil, err
		}

		return msgs, nil
	}

	log.Error("Unknown chain type: ", chain)

	return nil, nil
}

func (p *DefaultTxOutputProducer) getBridge(chain string) chainstypes.Bridge {
	bridge := p.bridges[chain]
	if bridge == nil {
		if libchain.IsETHBasedChain(chain) {
			p.bridges[chain] = chainseth.NewBridge(chain, p.signer, p.keeper)
		}
	}

	return p.bridges[chain]
}
