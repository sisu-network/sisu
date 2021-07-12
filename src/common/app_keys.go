package common

import (
	"errors"
	"os"

	keyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
)

type AppKeys struct {
	signerInfo keyring.Info
	kr         keyring.Keyring
	cfg        config.SisuConfig
}

func NewAppKeys(cfg config.SisuConfig) *AppKeys {
	return &AppKeys{
		cfg: cfg,
	}
}

func (ak *AppKeys) Init() {
	var err error
	utils.LogInfo("ak.cfg.KeyringBackend =", ak.cfg.KeyringBackend)
	utils.LogInfo("ak.cfg.Home =", ak.cfg.Home)

	ak.kr, err = keyring.New(sdk.KeyringServiceName(), ak.cfg.KeyringBackend, ak.cfg.Home, os.Stdin)
	if err != nil {
		panic(err)
	}

	infos, err := ak.kr.List()
	if len(infos) == 0 {
		utils.LogError()
		panic(errors.New(`Please create at least one account before running this node.
If this is a localhost network, run the gen file. If this is a testnet or
mainnet, generate account using "sisu keys" command.`))
	}

	// TODO: Use signer name for
	ak.signerInfo = infos[0]
	utils.LogDebug("signerInfo =", ak.signerInfo.GetName())
}

func (ak *AppKeys) GetSignerInfo() keyring.Info {
	return ak.signerInfo
}

func (ak *AppKeys) GetSignerAddress() sdk.AccAddress {
	return ak.signerInfo.GetAddress()
}

func (ak *AppKeys) GetKeyring() keyring.Keyring {
	return ak.kr
}
