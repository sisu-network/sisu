package dev

import (
	"encoding/hex"
	"fmt"
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
	liskConfig := helper.ReadLiskConfig(genesisFolder)
	fmt.Println("liskConfig.RPC = ", liskConfig.RPC)
	deyesChainCfg := config.Chain{Chain: liskConfig.Chain, Rpcs: []string{liskConfig.RPC}}

	client := deyeslisk.NewLiskClient(deyesChainCfg)
	mpcAddr := liskcrypto.GetAddressFromPublicKey(mpcPubKey)
	log.Verbose("Funding LSK for mpc address = ", mpcAddr)

	amount := uint64(20000000)
	moduleId := uint32(2)
	assetId := uint32(0)

	transferLisk(client, mnemonic, mpcAddr, amount, moduleId, assetId, liskConfig)
}

func transferLisk(client deyeslisk.Client, mnemonic, receiver string, amount uint64, moduleId uint32,
	assetId uint32, config helper.LiskConfig) {
	privateKey := liskcrypto.GetPrivateKeyFromSecret(mnemonic)
	faucet := liskcrypto.GetPublicKeyFromSecret(mnemonic)
	address := liskcrypto.GetAddressFromPublicKey(faucet)

	fmt.Println("Faucet address = ", address)

	senderAddress, err := hex.DecodeString(address)
	if err != nil {
		panic(err)
	}

	lisk32 := liskcrypto.AddressToLisk32(senderAddress)
	fmt.Println("Lisk Sender address = ", lisk32)
	acc, err := client.GetAccount(lisk32)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sender address = ", acc.Token)

	nonce, err := strconv.ParseUint(acc.Sequence.Nonce, 10, 64)
	if err != nil {
		panic(err)
	}

	recipientAddress, err := hex.DecodeString(receiver)
	if err != nil {
		panic(err)
	}

	fee := uint64(500000)
	data := ""
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
		SenderPublicKey: faucet,
	}
	txHash, err := proto.Marshal(txPb)
	if err != nil {
		panic(err)
	}

	signature := liskcrypto.SignWithNetwork(config.Network, txHash, privateKey)
	if err != nil {
		panic(err)
	}

	txPb.Signatures = [][]byte{signature}

	txHash, err = proto.Marshal(txPb)
	if err != nil {
		panic(err)
	}

	tx, err := client.CreateTransaction(hex.EncodeToString(txHash))
	if err != nil {
		panic(err)
	}

	log.Info("Lisk txHash = ", tx)
}
