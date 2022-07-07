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

type HandlerPauseContract struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
	mc     ManagerContainer
}

func NewHandlerPauseContract(mc ManagerContainer) *HandlerPauseContract {
	return &HandlerPauseContract{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
		mc:     mc,
	}
}

func (h *HandlerPauseContract) DeliverMsg(ctx sdk.Context, msg *types.PauseContractMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		data, err := newHandlerPauseResumeContract(h.mc).doPauseOrResume(ctx, msg.Data.Chain, msg.Data.Hash, true)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	} else {
		log.Verbose("HandlerPause: transaction has been processed")
	}

	return &sdk.Result{}, nil
}

////////////////////////////////
// HandlerResumeContract handles resuming contracts.
////////////////////////////////

type HandlerResumeContract struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
	mc     ManagerContainer
}

func NewHandlerResumeContract(mc ManagerContainer) *HandlerResumeContract {
	return &HandlerResumeContract{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
		mc:     mc,
	}
}

func (h *HandlerResumeContract) DeliverMsg(ctx sdk.Context, msg *types.ResumeContractMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		newHandlerPauseResumeContract(h.mc).doPauseOrResume(ctx, msg.Data.Chain, msg.Data.Hash, false)
		h.keeper.ProcessTxRecord(ctx, hash)
	} else {
		log.Verbose("HandlerResumeContract: transaction has been processed")
	}

	return &sdk.Result{}, nil
}

////////////////////////////////
// HandlerPauseResumeContract is a handler used for both pausing and resuming contracts.
////////////////////////////////

type handlerPauseResumeContract struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	partyManager     PartyManager
	dheartClient     tssclients.DheartClient
}

func newHandlerPauseResumeContract(mc ManagerContainer) *handlerPauseResumeContract {
	return &handlerPauseResumeContract{
		keeper:           mc.Keeper(),
		txOutputProducer: mc.TxOutProducer(),
		globalData:       mc.GlobalData(),
		partyManager:     mc.PartyManager(),
		dheartClient:     mc.DheartClient(),
	}
}

func (h *handlerPauseResumeContract) doPauseOrResume(ctx sdk.Context, chain, hash string, isPause bool) ([]byte, error) {
	// Only do pause/pause if we finished catching up.
	if h.globalData.IsCatchingUp() {
		log.Info("We are catching up with the network, exiting doPauseOrResume")
		return nil, nil
	}

	found := false
	for _, contract := range SupportedContracts {
		if contract.AbiHash == hash {
			found = true
			break
		}
	}

	if !found {
		err := fmt.Errorf("doPauseOrResume: contract with hash %s is not supported", hash)
		log.Error(err)
		return nil, err
	}

	var txOutMsg *types.TxOutWithSigner
	var err error
	if isPause {
		log.Info("Creating pause transaction...")
		txOutMsg, err = h.txOutputProducer.PauseContract(ctx, chain, hash)
	} else {
		log.Info("Creating resume transaction...")
		txOutMsg, err = h.txOutputProducer.ResumeContract(ctx, chain, hash)
	}

	if err != nil {
		log.Error("cannot get txOut for pausing contract, err = ", err)
		return nil, nil
	}

	// Save this to KVStore
	h.keeper.SaveTxOut(ctx, txOutMsg.Data)

	// Sends to dheart for signing.
	h.signTx(ctx, txOutMsg.Data)

	return nil, nil
}

// signTx sends a TxOut to dheart for TSS signing.
func (h *handlerPauseResumeContract) signTx(ctx sdk.Context, tx *types.TxOut) {
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

func (h *handlerPauseResumeContract) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
