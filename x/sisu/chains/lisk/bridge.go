package lisk

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ltypes "github.com/sisu-network/deyes/chains/lisk/types"
	eyestypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/chains/lisk/crypto"
	"math/big"

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
	responseTx, err := b.buildTransferInResponse(ctx, transfer.ToRecipient, transfer.Amount, "", moduleId, assetId)
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
	ctx sdk.Context,
	recipient string,
	amount string,
	data string,
	moduleId uint32,
	assetId uint32,
) (*types.TxResponse, error) {
	//config := deyesConfig.Chain{Chain: b.chain, Rpcs: []string{b.config.Lisk.RPC}}
	//client := NewLiskRPC(config)

	//senderAddress, err := hex.DecodeString()
	//if err != nil {
	//	return nil, err
	//}

	//senderLisk32 := crypto.AddressToLisk32(senderAddress)
	//acc, err := client.GetAccount(senderLisk32)
	//if err != nil {
	//	return nil, err
	//}

	//nonce, err := strconv.ParseUint(acc.Sequence.Nonce, 10, 32)
	//if err != nil {
	//	return nil, err
	//}

	recipientAddress, err := hex.DecodeString(recipient)
	if err != nil {
		return nil, err
	}

	fee := uint64(500000)
	amountInt := new(big.Int)
	amountOut, ok := amountInt.SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("SetString: error")
	}
	commissionRate := b.keeper.GetParams(ctx).CommissionRate
	if commissionRate < 0 || commissionRate > 10_000 {
		return nil, fmt.Errorf("commission rate is invalid, rate = %d", commissionRate)
	}
	// Subtract commission rate
	amountOut = utils.SubtractCommissionRate(amountOut, commissionRate)
	amountOut.Sub(amountOut, new(big.Int).SetUint64(fee))
	amountOutUint64 := amountOut.Uint64()
	assetPb := &ltypes.AssetMessage{
		Amount:           &amountOutUint64,
		RecipientAddress: recipientAddress,
		Data:             &data,
	}

	asset, err := proto.Marshal(assetPb)
	pubKey := crypto.GetPublicKeyFromSecret(b.config.Sisu.KeyringPassphrase)

	txPb := &ltypes.TransactionMessage{
		ModuleID: &moduleId,
		AssetID:  &assetId,
		Fee:      &fee,
		Asset:    asset,
		//Nonce:           &nonce,
		SenderPublicKey: pubKey,
	}
	txHash, err := proto.Marshal(txPb)
	if err != nil {
		return nil, err
	}

	txHash, err = proto.Marshal(txPb)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{
		OutChain: b.chain,
		RawBytes: txHash,
	}, nil
}

func (b *defaultBridge) ParseIncomingTx(ctx sdk.Context, chain string, tx *eyestypes.Tx) ([]*types.Transfer, error) {
	log.Verbose("Parsing lisk incoming tx...")
	ret := make([]*types.Transfer, 0)
	outerTx := ltypes.TransactionMessage{}
	proto.Unmarshal(tx.Serialized, &outerTx)

	if outerTx.AssetID == nil || outerTx.ModuleID == nil {
		return nil, fmt.Errorf("invalid outerTx")
	}

	asset := ltypes.AssetMessage{}
	proto.Unmarshal(outerTx.Asset, &asset)

	ret = append(ret, &types.Transfer{
		FromChain:   chain,
		FromHash:    hex.EncodeToString(outerTx.GetSignatures()[0]),
		Token:       "LSK",
		Amount:      strconv.FormatUint(asset.GetAmount(), 10),
		ToChain:     "",
		ToRecipient: hex.EncodeToString(asset.GetRecipientAddress()),
	})

	return ret, nil
}
