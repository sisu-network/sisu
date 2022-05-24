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

type HandlerContractChangeOwnership struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
	mc     ManagerContainer
}

func NewHandlerContractChangeOwnership(mc ManagerContainer) *HandlerContractChangeOwnership {
	return &HandlerContractChangeOwnership{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
		mc:     mc,
	}
}

func (h *HandlerContractChangeOwnership) DeliverMsg(ctx sdk.Context, msg *types.ChangeOwnershipContractMsg) (*sdk.Result, error) {
	process, hash, err := h.pmm.PreProcessingMsg(ctx, msg)
	if err != nil {
		return &sdk.Result{}, err
	}

	if !process {
		log.Verbose("HandlerContractChangeOwnership: didn't not reach consensus or transaction has been processed")
		return &sdk.Result{}, nil
	}

	data, err := newHandlerContractChangeOwnership(h.mc).doChangeOwner(ctx, msg.Data.Chain, msg.Data.Hash, msg.Data.NewOwner)
	h.keeper.ProcessTxRecord(ctx, hash)

	return &sdk.Result{Data: data}, err
}

type handlerContractChangeOwnership struct {
	keeper           keeper.Keeper
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	partyManager     PartyManager
	dheartClient     tssclients.DheartClient
}

func newHandlerContractChangeOwnership(mc ManagerContainer) *handlerContractChangeOwnership {
	return &handlerContractChangeOwnership{
		keeper:           mc.Keeper(),
		txOutputProducer: mc.TxOutProducer(),
		globalData:       mc.GlobalData(),
		partyManager:     mc.PartyManager(),
		dheartClient:     mc.DheartClient(),
	}
}

func (h *handlerContractChangeOwnership) doChangeOwner(ctx sdk.Context, chain, hash, newOwner string) ([]byte, error) {
	// Only do pause/pause if we finished catching up.
	if h.globalData.IsCatchingUp() {
		log.Info("We are catching up with the network, exiting doPauseOrResume")
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
		err := fmt.Errorf("doChangeOwner: contract with hash %s is not supported", hash)
		log.Error(err)
		return nil, err
	}

	msg, err := h.txOutputProducer.ContractChangeOwnership(ctx, chain, hash, newOwner)
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
func (h *handlerContractChangeOwnership) signTx(ctx sdk.Context, tx *types.TxOut) {
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

func (h *handlerContractChangeOwnership) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
