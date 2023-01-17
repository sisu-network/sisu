package dev

import (
	"encoding/hex"
	"strconv"

	deyeslisk "github.com/sisu-network/deyes/chains/lisk"
	liskcrypto "github.com/sisu-network/deyes/chains/lisk/crypto"
	lisktypes "github.com/sisu-network/deyes/chains/lisk/types"
	"github.com/sisu-network/deyes/config"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"google.golang.org/protobuf/proto"
)

func (c *fundAccountCmd) fundLisk(genesisFolder, mnemonic string, mpcPubKey []byte) {
	amount := uint64(1 * 100_000_000)
	transferLisk(genesisFolder, mnemonic, mpcPubKey, amount, "")
}

func transferLisk(genesisFolder, mnemonic string, mpcPubKey []byte, amount uint64, data string) {
	liskConfig := helper.ReadLiskConfig(genesisFolder)
	log.Verbosef("Use url %s for chain %s", liskConfig.RPC, liskConfig.Chain)
	deyesChainCfg := config.Chain{Chain: liskConfig.Chain, Rpcs: []string{liskConfig.RPC}}

	client := deyeslisk.NewLiskClient(deyesChainCfg)
	mpcAddr := liskcrypto.GetAddressFromPublicKey(mpcPubKey)
	log.Verbose("Funding LSK for mpc address = ", mpcAddr)

	receiver := mpcAddr

	moduleId := uint32(2)
	assetId := uint32(0)

	privateKey := liskcrypto.GetPrivateKeyFromSecret(mnemonic)
	faucetPubKey := liskcrypto.GetPublicKeyFromSecret(mnemonic)

	lisk32 := liskcrypto.GetLisk32AddressFromPublickey(faucetPubKey)
	acc, err := client.GetAccount(lisk32)
	if err != nil {
		panic(err)
	}

	nonce, err := strconv.ParseUint(acc.Sequence.Nonce, 10, 64)
	if err != nil {
		panic(err)
	}

	recipientAddress, err := hex.DecodeString(receiver)
	if err != nil {
		panic(err)
	}

	fee := uint64(500_000)
	assetPb := &lisktypes.AssetMessage{
		Amount:           &amount,
		RecipientAddress: recipientAddress,
		Data:             &data,
	}

	asset, err := proto.Marshal(assetPb)
	txPb := &lisktypes.TransactionMessage{
		ModuleID:        &moduleId,
		AssetID:         &assetId,
		Fee:             &fee,
		Asset:           asset,
		Nonce:           &nonce,
		SenderPublicKey: faucetPubKey,
	}
	txHash, err := proto.Marshal(txPb)
	if err != nil {
		panic(err)
	}

	signature := liskcrypto.SignWithNetwork(liskConfig.Network, txHash, privateKey)
	if err != nil {
		panic(err)
	}

	txPb.Signatures = [][]byte{signature}

	txHash, err = proto.Marshal(txPb)
	if err != nil {
		panic(err)
	}

	log.Infof("Funding Sisu from account %s to account %s= ", lisk32,
		liskcrypto.GetLisk32AddressFromPublickey(mpcPubKey))

	tx, err := client.CreateTransaction(hex.EncodeToString(txHash))
	if err != nil {
		panic(err)
	}

	log.Info("Lisk txHash = ", tx)
}
