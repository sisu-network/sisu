package dev

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	ltypes "github.com/sisu-network/deyes/chains/lisk/types"
	"github.com/sisu-network/deyes/config"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/helper"
	"github.com/sisu-network/sisu/x/sisu/chains/lisk"
	"github.com/sisu-network/sisu/x/sisu/chains/lisk/crypto"
	"strconv"
)

func (c *fundAccountCmd) fundLisk(genesisFolder, mnemonic string, mpcPubKey []byte) {
	liskConfig := helper.ReadLiskConfig(genesisFolder)
	configFormatted := config.Chain{Chain: liskConfig.Chain, Rpcs: []string{liskConfig.RPC}}
	client := lisk.NewLiskRPC(configFormatted)
	mpcAddr := crypto.GetAddressFromPublicKey(mpcPubKey)
	log.Verbose("Funding LSK for mpc address = ", mpcAddr)
	amount := uint64(20000000)
	moduleId := uint32(2)
	assetId := uint32(0)
	transferLisk(client, mnemonic, mpcAddr, amount, moduleId, assetId, liskConfig)
}

func transferLisk(client lisk.LiskRPC, mnemonic, receiver string, amount uint64, moduleId uint32, assetId uint32, config helper.LiskConfig) {
	privateKey := crypto.GetPrivateKeyFromSecret(mnemonic)
	faucet := crypto.GetPublicKeyFromSecret(mnemonic)
	address := crypto.GetAddressFromPublicKey(faucet)
	senderAddress, err := hex.DecodeString(address)
	if err != nil {
		panic(err)
	}

	lisk32 := crypto.AddressToLisk32(senderAddress)
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

	fee := uint64(500000)

	data := "fund sisu"
	assetPb := &ltypes.AssetMessage{
		Amount:           &amount,
		RecipientAddress: recipientAddress,
		Data:             &data,
	}

	asset, err := proto.Marshal(assetPb)
	txPb := &ltypes.TransactionMessage{
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

	signature := sign(config.Network, txHash, privateKey)
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

func sign(network string, txBytes []byte, privateKey []byte) []byte {
	dst := new(bytes.Buffer)
	//First byte is the network info
	networkBytes, err := hex.DecodeString(network)
	if err != nil {
		panic(err)
	}
	binary.Write(dst, binary.LittleEndian, networkBytes)

	// Append the transaction ModuleID
	binary.Write(dst, binary.LittleEndian, txBytes)

	return crypto.SignMessageWithPrivateKey(string(dst.Bytes()), privateKey)
}
