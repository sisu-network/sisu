package common

import (
	"encoding/hex"
	"errors"
	"os"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/sisu-network/lib/log"

	keyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
)

//go:generate mockgen -source=common/app_keys.go -destination=./tests/mock/common/app_keys.go -package=mock

// Make sure struct implement interface at compile-time
var _ AppKeys = (*DefaultAppKeys)(nil)

type AppKeys interface {
	Init()
	GetSignerInfo() keyring.Info
	GetSignerAddress() sdk.AccAddress
	GetKeyring() keyring.Keyring
	GetEncryptedPrivKey() ([]byte, error)
	GetAesEncrypted(msg []byte) ([]byte, error)
}

type DefaultAppKeys struct {
	signerInfo keyring.Info
	kr         keyring.Keyring
	cfg        config.SisuConfig
	privateKey ctypes.PrivKey
	aesKey     []byte
}

func NewAppKeys(cfg config.SisuConfig) *DefaultAppKeys {
	return &DefaultAppKeys{
		cfg: cfg,
	}
}

func (ak *DefaultAppKeys) Init() {
	var err error
	log.Info("ak.cfg.KeyringBackend =", ak.cfg.KeyringBackend)
	log.Info("ak.cfg.Home =", ak.cfg.Dir)

	ak.kr, err = keyring.New(sdk.KeyringServiceName(), ak.cfg.KeyringBackend, ak.cfg.Dir, os.Stdin)
	if err != nil {
		panic(err)
	}

	infos, err := ak.kr.List()
	if len(infos) == 0 {
		log.Error()
		panic(errors.New(`Please create at least one account before running this node.
If this is a localhost network, run the gen file.
If this is a testnet or mainnet, generate account using "sisu keys" command.`))
	}

	// TODO: Use signer name for
	ak.signerInfo = infos[0]
	// Set the private key from keyring
	ak.setPrivateKey()

	// Set the AES Key
	ak.aesKey, err = hex.DecodeString(os.Getenv("AES_KEY_HEX"))
	if err != nil {
		panic("invalid aes key")
	}
}

func (ak *DefaultAppKeys) setPrivateKey() {
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

func (ak *DefaultAppKeys) GetSignerInfo() keyring.Info {
	return ak.signerInfo
}

func (ak *DefaultAppKeys) GetSignerAddress() sdk.AccAddress {
	return ak.signerInfo.GetAddress()
}

func (ak *DefaultAppKeys) GetKeyring() keyring.Keyring {
	return ak.kr
}

func (ak *DefaultAppKeys) GetEncryptedPrivKey() ([]byte, error) {
	bz := ak.privateKey.Bytes()
	return utils.AESDEncrypt(bz, ak.aesKey)
}

func (ak *DefaultAppKeys) GetAesEncrypted(msg []byte) ([]byte, error) {
	return utils.AESDEncrypt(msg, ak.aesKey)
}
