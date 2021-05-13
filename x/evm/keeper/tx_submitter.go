package keeper

import (
	"bytes"
	"io"
	"os"

	dcore "github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/sisu/app/params"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/evm/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	defaultGasAdjustment = 1.0
	defaultGasLimit      = 300000
)

type TxSubmitter struct {
	sisuHome  string
	kr        keyring.Keyring
	clientCtx client.Context
	factory   tx.Factory
}

var (
	nodeAddress = "http://0.0.0.0:26657"
	// TODO: Use correct chain id
	chainId = "chain-Gbme39"
)

func NewTxSubmitter(sisuHome string, keyRingBackend string) *TxSubmitter {
	// TODO: Fix this
	kb, err := keyring.New(sdk.KeyringServiceName(), keyRingBackend, sisuHome, os.Stdin)
	if err != nil {
		utils.LogError("Cannot create keyring")
		return nil
	}

	t := &TxSubmitter{
		kr: kb,
	}

	infos, err := kb.List()

	t.clientCtx, err = t.buildClientCtx(infos[0].GetName())
	t.factory = newFactory(t.clientCtx)

	if err != nil {
		return nil
	}

	return t
}

func (t *TxSubmitter) onTxSubmitted(ethTx *dcore.Transaction) {
	go func() {
		js, err := ethTx.MarshalJSON()
		if err != nil {
			return
		}

		msg := types.NewMsgEthTx(t.clientCtx.GetFromAddress().String(), js)
		if err := tx.BroadcastTx(t.clientCtx, t.factory, msg); err != nil {
			utils.LogError("Cannot broadcast transaction", err)
			return
		}
	}()
}

func (t *TxSubmitter) buildClientCtx(accountName string) (client.Context, error) {
	info, err := t.kr.Key(accountName)
	if err != nil {
		return client.Context{}, err
	}

	client, err := rpchttp.New(nodeAddress, "/websocket")
	clientCtx := NewClientCtx(t.kr, client, &bytes.Buffer{}, t.sisuHome)

	return clientCtx.
		WithFromName(accountName).
		WithFromAddress(info.GetAddress()), nil
}

func NewClientCtx(kr keyring.Keyring, c *rpchttp.HTTP, out io.Writer, home string) client.Context {
	encodingConfig := params.MakeEncodingConfig()
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sdk.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	staking.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return client.Context{}.
		WithChainID(chainId).
		WithKeyring(kr).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(out).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(home).
		WithClient(c).
		WithSkipConfirmation(true)
}

func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(defaultGasLimit).
		WithGasAdjustment(defaultGasAdjustment).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig)
}
