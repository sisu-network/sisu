package common

import (
	"encoding/hex"
	"errors"
	"os"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"

	keyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
)

type AppKeys struct {
	signerInfo keyring.Info
	kr         keyring.Keyring
	cfg        config.SisuConfig
	privateKey ctypes.PrivKey
	aesKey     []byte
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

	// Set the private key from keyring
	ak.setPrivateKey()

	// Set the AES Key
	ak.aesKey, err = hex.DecodeString(os.Getenv("AES_KEY_HEX"))
	if err != nil {
		panic("invalid aes key")
	}
}

func (ak *AppKeys) setPrivateKey() {
	keyType := ak.signerInfo.GetPubKey().Type()
	unsafe := keyring.NewUnsafe(ak.kr)
	hexKey, err := unsafe.UnsafeExportPrivKeyHex(ak.signerInfo.GetName())
	if err != nil {
		panic(err)
	}

	bz, err := hex.DecodeString(hexKey)
	if err != nil {
		panic(err)
	}

	if keyType == "secp256k1" {
		ak.privateKey = &secp256k1.PrivKey{Key: bz}
	} else {
		panic("unsupported key type")
	}
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

func (ak *AppKeys) GetEncryptedPrivKey() ([]byte, error) {
	bz := ak.privateKey.Bytes()
	return utils.AESDEncrypt(bz, ak.aesKey)
}
