package lisk

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	libchain "github.com/sisu-network/lib/chain"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/near/borsh-go"
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

	nonce, err := b.deyesClient.GetNonce(b.chain, mpcAddr)
	if err != nil {
		return nil, err
	}

	recipient, err := hex.DecodeString(transfer.ToRecipient)
	if err != nil {
		return nil, err
	}

	amountInt, ok := new(big.Int).SetString(transfer.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("Invalid transfer amount %s", transfer.Amount)
	}

	commissionRate := b.keeper.GetParams(ctx).CommissionRate
	if commissionRate < 0 || commissionRate > 10_000 {
		return nil, fmt.Errorf("commission rate is invalid, rate = %d", commissionRate)
	}

	fee := uint64(500000)
	amountOut := new(big.Int).Set(amountInt)

	// Subtract commission rate
	amountOut = utils.SubtractCommissionRate(amountOut, commissionRate)
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

	moduleId := uint32(2)
	assetId := uint32(0)
	txPb := &lisktypes.TransactionMessage{
		ModuleID:        &moduleId,
		AssetID:         &assetId,
		Fee:             &fee,
		Asset:           asset,
		Nonce:           &nonce,
		SenderPublicKey: []byte(mpcAddr), // TODO: check if this is correct
	}

	bz, err := proto.Marshal(txPb)
	if err != nil {
		return nil, err
	}

	networkId := lisktypes.NetworkId[transfer.ToChain]
	if len(networkId) == 0 {
		return nil, fmt.Errorf("cannot find lisk network id for chain %s", transfer.ToChain)
	}

	// The signing bytes are the combination of network id and the serialization of the transaciton.
	signBuf := new(bytes.Buffer)
	//First byte is the network info
	networkBytes, _ := hex.DecodeString(networkId)
	binary.Write(signBuf, binary.LittleEndian, networkBytes)

	// Append the transaction ModuleID
	binary.Write(signBuf, binary.LittleEndian, bz)

	outMsg := types.NewTxOutMsg(
		b.signer,
		types.TxOutType_TRANSFER_OUT,
		&types.TxOutContent{
			OutChain: b.chain,
			OutHash:  hex.EncodeToString(bz),
			OutBytes: signBuf.Bytes(),
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
	err = borsh.Deserialize(payload, bz)
	if err != nil {
		return nil, err
	}

	dstChain := libchain.GetChainNameFromInt(big.NewInt(int64(payload.ChainId)))
	if len(dstChain) == 0 {
		return nil, fmt.Errorf("Invalid destination chain %s", dstChain)
	}

	token := b.keeper.GetToken(ctx, payload.Token)
	if token == nil {
		return nil, fmt.Errorf("Invalid token %s", payload.Token)
	}

	amount, err := token.ConvertAmountToSisuAmount(b.chain, big.NewInt(int64(payload.Amount)))
	if err != nil {
		return nil, err
	}

	fmt.Println("Parsed amount = ", amount)

	return []*types.TransferDetails{
		{
			Id:        fmt.Sprintf("%s__%s", chain, tx.Id),
			FromChain: chain,
			ToChain:   dstChain,
			Token:     payload.Token,
			Amount:    amount.String(),
		},
	}, nil
}
