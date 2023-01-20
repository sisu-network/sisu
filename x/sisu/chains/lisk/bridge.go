package lisk

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	liskcrypto "github.com/sisu-network/deyes/chains/lisk/crypto"
	lisktypes "github.com/sisu-network/deyes/chains/lisk/types"
	"github.com/sisu-network/sisu/utils"
	chainstypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"google.golang.org/protobuf/proto"
)

type bridge struct {
	signer      string
	chain       string
	deyesClient external.DeyesClient
	keeper      keeper.Keeper
}

func NewBridge(chain string, signer string, keeper keeper.Keeper, deyesClient external.DeyesClient,
) chainstypes.Bridge {
	return &bridge{
		signer:      signer,
		chain:       chain,
		deyesClient: deyesClient,
		keeper:      keeper,
	}
}

func (b *bridge) ProcessTransfers(ctx sdk.Context, transfers []*types.TransferDetails) (
	[]*types.TxOutMsg, error) {
	// Lisk only accept 1 recipient
	if len(transfers) != 1 {
		return nil, fmt.Errorf("Lisk only accept 1 recipient")
	}

	transfer := transfers[0]
	mpcAddr := b.keeper.GetMpcAddress(ctx, b.chain)
	mpcPubkey := b.keeper.GetMpcPublicKey(ctx, b.chain)

	nonce, err := b.deyesClient.GetNonce(b.chain, mpcAddr)
	if err != nil {
		return nil, err
	}

	recipient, err := liskcrypto.Lisk32AddressToPublicAddress(transfer.ToRecipient)
	if err != nil {
		return nil, err
	}

	amountInt, ok := new(big.Int).SetString(transfer.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("Invalid transfer amount %s", transfer.Amount)
	}

	token := b.keeper.GetToken(ctx, transfer.Token)

	fee := uint64(500_000)
	amountOut := new(big.Int).Set(amountInt)
	// Conver the amount from Sisu amount (18 decimals to Lisk 8 decimals)
	amountOut, err = token.SisuAmountToChainAmount(transfer.ToChain, amountOut)
	if err != nil {
		return nil, err
	}

	// Subtract commission rate
	commissionRate := b.keeper.GetParams(ctx).CommissionRate
	if commissionRate < 0 || commissionRate > 10_000 {
		return nil, fmt.Errorf("commission rate is invalid, rate = %d", commissionRate)
	}
	amountOut = utils.SubtractCommissionRate(amountOut, commissionRate)

	// Subtract transaction fee
	amountOut.Sub(amountOut, new(big.Int).SetUint64(fee))

	amountOutUint64 := amountOut.Uint64()
	data := ""
	assetPb := &lisktypes.AssetMessage{
		Amount:           &amountOutUint64,
		RecipientAddress: recipient,
		Data:             &data,
	}
	asset, err := proto.Marshal(assetPb)
	if err != nil {
		return nil, err
	}

	log.Verbosef("Lisk transfer details, amount out = %d, recipient = %s, nonce = %d",
		amountOut.Uint64(), hex.EncodeToString(recipient), nonce)

	moduleId := uint32(2)
	assetId := uint32(0)
	tx := &lisktypes.TransactionMessage{
		ModuleID:        &moduleId,
		AssetID:         &assetId,
		Fee:             &fee,
		Asset:           asset,
		Nonce:           &nonce,
		SenderPublicKey: mpcPubkey,
	}

	bz, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}

	hash, err := tx.GetHash()
	if err != nil {
		return nil, err
	}

	outMsg := types.NewTxOutMsg(
		b.signer,
		types.TxOutType_TRANSFER_OUT,
		&types.TxOutContent{
			OutChain: b.chain,
			OutHash:  hash,
			OutBytes: bz,
		},
		&types.TxOutInput{
			TransferIds: []string{transfer.Id},
		},
	)

	return []*types.TxOutMsg{outMsg}, nil
}

func (b *bridge) ParseIncomingTx(ctx sdk.Context, chain string, serialized []byte) (
	[]*types.TransferDetails, error) {
	tx := &lisktypes.Transaction{}
	err := json.Unmarshal(serialized, tx)
	if err != nil {
		return nil, err
	}

	bz, err := base64.StdEncoding.DecodeString(tx.Asset.Data)
	if err != nil {
		return nil, err
	}

	payload := &lisktypes.TransferData{}
	err = proto.Unmarshal(bz, payload)
	if err != nil {
		return nil, err
	}

	dstChain := libchain.GetChainNameFromInt(big.NewInt(int64(*payload.ChainId)))
	if len(dstChain) == 0 {
		return nil, fmt.Errorf("Invalid destination chain %s", dstChain)
	}

	tokenStr := "LSK"
	if payload.Token != nil {
		tokenStr = *payload.Token
	}
	token := b.keeper.GetToken(ctx, tokenStr)
	if token == nil {
		return nil, fmt.Errorf("Invalid token %s", tokenStr)
	}

	amount, err := token.ChainAmountToSisuAmount(b.chain, big.NewInt(int64(*payload.Amount)))
	if err != nil {
		return nil, err
	}

	var recipient string
	switch {
	case libchain.IsETHBasedChain(dstChain):
		recipient = fmt.Sprintf("0x%s", hex.EncodeToString(payload.Recipient))
	default:
		return nil, fmt.Errorf("unsupported lisk payload address for chain %s, recipient hex = %s",
			dstChain, hex.EncodeToString(payload.Recipient))
	}

	return []*types.TransferDetails{
		{
			Id:          fmt.Sprintf("%s__%s", chain, tx.Id),
			FromChain:   chain,
			ToChain:     dstChain,
			Token:       tokenStr,
			ToRecipient: recipient,
			Amount:      amount.String(),
		},
	}, nil
}
