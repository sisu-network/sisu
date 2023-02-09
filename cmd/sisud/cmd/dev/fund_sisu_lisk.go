package dev

import (
	"crypto/sha256"
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
	amount := uint64(3 * 100_000_000)
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
	log.Verbosef("Lisk32 of the faucet = %s", lisk32)
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
	tx := &lisktypes.TransactionMessage{
		ModuleID:        &moduleId,
		AssetID:         &assetId,
		Fee:             &fee,
		Asset:           asset,
		Nonce:           &nonce,
		SenderPublicKey: faucetPubKey,
	}
	bz, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	bytesToSign, err := liskcrypto.GetSigningBytes(lisktypes.NetworkId[liskConfig.Chain], bz)
	if err != nil {
		log.Errorf("Failed to get lisk bytes to sign, err = %s", err)
		return
	}

	signature := liskcrypto.SignMessage(bytesToSign, privateKey)
	tx.Signatures = [][]byte{signature}
	signedBz, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(signedBz)
	log.Verbosef("Calculated hash = %s", hex.EncodeToString(hash[:]))
	log.Infof("Funding Sisu from account %s to account %s", lisk32,
		liskcrypto.GetLisk32AddressFromPublickey(mpcPubKey))

	txHash, err := client.CreateTransaction(hex.EncodeToString(signedBz))
	if err != nil {
		panic(err)
	}

	log.Info("Lisk txHash = ", txHash)
}
