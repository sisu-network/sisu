package sisu

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/ethereum/go-ethereum/core/types"
	hTypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerContractSetLiquidityAddress struct {
	pmm    PostedMessageManager
	mc     ManagerContainer
	keeper keeper.Keeper
}

func NewHandlerContractSetLiquidityAddress(mc ManagerContainer) *HandlerContractSetLiquidityAddress {
	return &HandlerContractSetLiquidityAddress{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
		mc:     mc,
	}
}

func (h *HandlerContractSetLiquidityAddress) DeliverMsg(ctx sdk.Context, msg *types.ChangeLiquidPoolAddressMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		data, err := newHandlerContractSetLiquidityAddress(h.mc).doSetLiquidityAddress(ctx, msg.Data.Chain, msg.Data.Hash, msg.Data.NewLiquidAddress)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	} else {
		log.Verbose("HandlerContractSetLiquidityAddress: didn't not reach consensus or transaction has been processed")
	}

	return &sdk.Result{}, nil
}

type handlerContractSetLiquidityAddress struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	partyManager     PartyManager
	dheartClient     tssclients.DheartClient
}

func newHandlerContractSetLiquidityAddress(mc ManagerContainer) *handlerContractSetLiquidityAddress {
	return &handlerContractSetLiquidityAddress{
		keeper:           mc.Keeper(),
		txOutputProducer: mc.TxOutProducer(),
		globalData:       mc.GlobalData(),
		partyManager:     mc.PartyManager(),
		dheartClient:     mc.DheartClient(),
	}
}

func (h *handlerContractSetLiquidityAddress) doSetLiquidityAddress(ctx sdk.Context, chain, hash, newLpAddress string) ([]byte, error) {
	// Only do set liquidity address if we finished catching up.
	if h.globalData.IsCatchingUp() {
		log.Info("We are catching up with the network, exiting setLiquidAddress")
		return nil, nil
	}

	found := false
	for _, contract := range SupportedContracts {
		if strings.EqualFold(strings.ToLower(contract.AbiHash), strings.ToLower(hash)) {
			found = true
			break
		}
	}

	if !found {
		err := fmt.Errorf("doSetLiquidityAddress: contract with hash %s is not supported", hash)
		log.Error(err)
		return nil, err
	}

	txOutMsg, err := h.txOutputProducer.ContractSetLiquidPoolAddress(ctx, chain, hash, newLpAddress)
	if err != nil {
		return nil, err
	}

	// Save this to KVStore
	h.keeper.SaveTxOut(ctx, txOutMsg.Data)

	// Sends to dheart for signing.
	h.signTx(ctx, txOutMsg.Data)

	return nil, nil
}

// TODO: duplicate code with pause/resume contract handler, fix it
// signTx sends a TxOut to dheart for TSS signing.
func (h *handlerContractSetLiquidityAddress) signTx(ctx sdk.Context, tx *types.TxOut) {
	ethTx := &etypes.Transaction{}
	if err := ethTx.UnmarshalBinary(tx.OutBytes); err != nil {
		log.Error("cannot unmarshal tx, err =", err)
	}

	signer := libchain.GetEthChainSigner(tx.OutChain)
	if signer == nil {
		err := fmt.Errorf("cannot find signer for chain %s", tx.OutChain)
		log.Error(err)
	}

	hash := signer.Hash(ethTx)

	// Send it to Dheart for signing.
	keysignReq := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_ECDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          h.getKeysignRequestId(tx.OutChain, ctx.BlockHeight(), tx.OutHash),
				OutChain:    tx.OutChain,
				OutHash:     tx.OutHash,
				BytesToSign: hash[:],
			},
		},
	}

	pubKeys := h.partyManager.GetActivePartyPubkeys()

	err := h.dheartClient.KeySign(keysignReq, pubKeys)

	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func (h *handlerContractSetLiquidityAddress) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
