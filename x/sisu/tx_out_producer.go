package sisu

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"

	chainscar "github.com/sisu-network/sisu/x/sisu/chains/cardano"
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
	deyesClient    external.DeyesClient

	bridges map[string]chainstypes.Bridge
}

type transferInData struct {
	token     ethcommon.Address
	recipient string
	amount    *big.Int
}

func NewTxOutputProducer(appKeys common.AppKeys, keeper keeper.Keeper,
	cardanoConfig config.CardanoConfig,
	deyesClient external.DeyesClient,
	txTracker TxTracker) TxOutputProducer {
	return &DefaultTxOutputProducer{
		signer:         appKeys.GetSignerAddress().String(),
		keeper:         keeper,
		txTracker:      txTracker,
		cardanoNetwork: cardanoConfig.GetCardanoNetwork(),
		deyesClient:    deyesClient,
		bridges:        make(map[string]chainstypes.Bridge),
	}
}

func (p *DefaultTxOutputProducer) GetTxOuts(ctx sdk.Context, chain string,
	transfers []*types.Transfer) ([]*types.TxOutMsg, error) {

	bridge := p.getBridge(chain)
	msgs, err := bridge.ProcessTransfers(ctx, transfers)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (p *DefaultTxOutputProducer) getBridge(chain string) chainstypes.Bridge {
	bridge := p.bridges[chain]
	if bridge == nil {
		if libchain.IsETHBasedChain(chain) {
			p.bridges[chain] = chainseth.NewBridge(chain, p.signer, p.keeper)
		}

		if libchain.IsCardanoChain(chain) {
			p.bridges[chain] = chainscar.NewBridge(chain, p.signer, p.keeper, p.deyesClient)
		}
	}

	return p.bridges[chain]
}
