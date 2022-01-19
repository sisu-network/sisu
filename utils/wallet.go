package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/sisu-network/sisu/utils/hdwallet"
)

const LOCALHOST_MNEMONIC = "draft attract behave allow rib raise puzzle frost neck curtain gentle bless letter parrot hold century diet budget paper fetch hat vanish wonder maximum"
const LOCAL_KEYSTORE_PASS = "kyrre"

var (
	localAccounts []accounts.Account
	localWallet   *hdwallet.Wallet
)

func genLocalhostAccounts() {
	var err error
	localWallet, err = hdwallet.NewFromMnemonic(LOCALHOST_MNEMONIC)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		account, err := localWallet.Derive(path, true)
		if err != nil {
			panic(err)
		}
		localAccounts = append(localAccounts, account)
	}
}

func GetLocalAccounts() []accounts.Account {
	return localAccounts
}

func GetLocalWallet() *hdwallet.Wallet {
	return localWallet
}
