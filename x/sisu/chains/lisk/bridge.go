package lisk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ltypes "github.com/sisu-network/deyes/chains/lisk/types"
	deyesConfig "github.com/sisu-network/deyes/config"
	eyestypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/chains/lisk/crypto"

	"github.com/golang/protobuf/proto"
	chaintypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"

	"strconv"
)

type defaultBridge struct {
	chain  string
	keeper keeper.Keeper
	signer string
	config config.Config
}

func NewBridge(chain, signer string, keeper keeper.Keeper, cfg config.Config) chaintypes.Bridge {
	return &defaultBridge{
		chain:  chain,
		keeper: keeper,
		signer: signer,
		config: cfg,
	}
}

func (b *defaultBridge) ProcessTransfers(ctx sdk.Context, transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
	moduleId := uint32(2)
	assetId := uint32(0)
	transfer := transfers[0]
	amount, err := strconv.ParseUint(transfer.Amount, 10, 64)
	if err != nil {
		return nil, err
	}
	responseTx, err := b.buildTransferInResponse("mnemonic", transfer.ToRecipient, amount, "", moduleId, assetId)
	if err != nil {
		log.Error("Failed to build lisk transfer in, err = ", err)
		return nil, err
	}
	outMsg := types.NewTxOutMsg(
		b.signer,
		types.TxOutType_TRANSFER_OUT,
		&types.TxOutContent{
			OutChain: b.chain,
			OutHash:  hex.EncodeToString(responseTx.RawBytes),
			OutBytes: responseTx.RawBytes,
		},
		&types.TxOutInput{
			TransferIds: []string{transfer.Id},
		},
	)
	return []*types.TxOutMsg{outMsg}, nil
}

func (b *defaultBridge) buildTransferInResponse(
	mnemonic string,
	recipient string,
	amount uint64,
	data string,
	moduleId uint32,
	assetId uint32,
) (*types.TxResponse, error) {
	config := deyesConfig.Chain{Chain: b.chain, Rpcs: []string{b.config.Lisk.RPC}}
	client := NewLiskRPC(config)
	faucet := crypto.GetPublicKeyFromSecret(mnemonic)
	address := crypto.GetAddressFromPublicKey(faucet)

	senderAddress, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	senderLisk32 := crypto.AddressToLisk32(senderAddress)
	acc, err := client.GetAccount(senderLisk32)
	if err != nil {
		return nil, err
	}

	nonce, err := strconv.ParseUint(acc.Sequence.Nonce, 10, 32)
	if err != nil {
		return nil, err
	}

	recipientAddress, err := hex.DecodeString(recipient)
	if err != nil {
		return nil, err
	}

	assetPb := &ltypes.AssetMessage{
		Amount:           &amount,
		RecipientAddress: recipientAddress,
		Data:             &data,
	}

	asset, err := proto.Marshal(assetPb)
	pubKey := crypto.GetPublicKeyFromSecret(b.config.Sisu.KeyringPassphrase)
	privateKey := crypto.GetPrivateKeyFromSecret(b.config.Sisu.KeyringPassphrase)
	fee := uint64(1000000)
	txPb := &ltypes.TransactionMessage{
		ModuleID:        &moduleId,
		AssetID:         &assetId,
		Fee:             &fee,
		Asset:           asset,
		Nonce:           &nonce,
		SenderPublicKey: pubKey,
	}
	txHash, err := proto.Marshal(txPb)
	if err != nil {
		return nil, err
	}

	signature, err := sign(b.config.Lisk.Network["testnet"], txHash, privateKey)
	if err != nil {
		return nil, err
	}
	txPb.Signatures = [][]byte{signature}

	txHash, err = proto.Marshal(txPb)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{
		OutChain: b.chain,
		RawBytes: txHash,
	}, nil
}

func (b *defaultBridge) ParseIncomginTx(ctx sdk.Context, chain string, tx *eyestypes.Tx) ([]*types.Transfer, error) {

	return nil, nil
}

func sign(network string, txBytes []byte, privateKey []byte) ([]byte, error) {
	dst := new(bytes.Buffer)
	//First byte is the network info
	networkBytes, _ := hex.DecodeString(network)
	binary.Write(dst, binary.LittleEndian, networkBytes)

	// Append the transaction ModuleID
	binary.Write(dst, binary.LittleEndian, txBytes)

	return crypto.SignMessageWithPrivateKey(string(dst.Bytes()), privateKey), nil
}
