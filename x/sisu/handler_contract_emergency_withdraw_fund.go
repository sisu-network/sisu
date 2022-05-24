package sisu

import (
	"fmt"
	"strconv"

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

type HandlerContractLiquidityWithdrawFund struct {
	pmm              PostedMessageManager
	mc               ManagerContainer
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	partyManager     PartyManager
	dheartClient     tssclients.DheartClient
}

func NewHandlerContractLiquidityWithdrawFund(mc ManagerContainer) *HandlerContractLiquidityWithdrawFund {
	return &HandlerContractLiquidityWithdrawFund{
		keeper:           mc.Keeper(),
		pmm:              mc.PostedMessageManager(),
		mc:               mc,
		txOutputProducer: mc.TxOutProducer(),
		globalData:       mc.GlobalData(),
		partyManager:     mc.PartyManager(),
		dheartClient:     mc.DheartClient(),
	}
}

func (h *HandlerContractLiquidityWithdrawFund) DeliverMsg(ctx sdk.Context, msg *types.LiquidityWithdrawFundMsg) (*sdk.Result, error) {
	if process, hash, err := h.pmm.ProcessMsg(ctx, msg); process {
		if err != nil {
			return &sdk.Result{}, err
		}
	
		data, err := h.doWithdrawFund(ctx, msg.Data.Chain, msg.Data.Hash, msg.Data.TokenAddresses, msg.Data.NewOwner)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	} else {
		log.Verbose("HandlerContractLiquidityWithdrawFund: didn't not reach consensus or transaction has been processed")
	}

	return &sdk.Result{}, nil
}

func (h *HandlerContractLiquidityWithdrawFund) doWithdrawFund(ctx sdk.Context, chain, hash string,
	tokens []string, newOwner string) ([]byte, error) {
	if h.globalData.IsCatchingUp() {
		log.Info("We are catching up with the network, exiting doWithdrawFund")
		return nil, nil
	}

	msg, err := h.txOutputProducer.ContractEmergencyWithdrawFund(ctx, chain, hash, tokens, newOwner)
	if err != nil {
		return nil, err
	}

	// Save this to KVStore
	h.keeper.SaveTxOut(ctx, msg.Data)

	// Sends to dheart for signing.
	h.signTx(ctx, msg.Data)

	return nil, nil
}

// TODO: duplicate code with pause/resume contract handler, fix it
// signTx sends a TxOut to dheart for TSS signing.
func (h *HandlerContractLiquidityWithdrawFund) signTx(ctx sdk.Context, tx *types.TxOut) {
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
				InChain:     tx.InChain,
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

func (h *HandlerContractLiquidityWithdrawFund) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
