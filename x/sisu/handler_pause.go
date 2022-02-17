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

type HandlerPause struct {
	pmm              PostedMessageManager
	publicDb         keeper.Storage
	txOutputProducer TxOutputProducer
	globalData       common.GlobalData
	partyManager     PartyManager
	dheartClient     tssclients.DheartClient
}

func NewHandlerPauseContract(mc ManagerContainer) *HandlerPause {
	return &HandlerPause{
		publicDb:         mc.PublicDb(),
		pmm:              mc.PostedMessageManager(),
		txOutputProducer: mc.TxOutProducer(),
		globalData:       mc.GlobalData(),
		partyManager:     mc.PartyManager(),
		dheartClient:     mc.DheartClient(),
	}
}

func (h *HandlerPause) DeliverMsg(ctx sdk.Context, msg *types.PauseContractMsg) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, msg); process {
		h.doPause(ctx, msg.Data)
		h.publicDb.ProcessTxRecord(hash)
	}

	return nil, nil
}

func (h *HandlerPause) doPause(ctx sdk.Context, pauseContract *types.PauseContract) ([]byte, error) {
	// Only do pause if we are catching up.
	if h.globalData.IsCatchingUp() {
		log.Info("We are catching up with the network, exiting doPause")
		return nil, nil
	}

	log.Info("Pausing contract...")

	if len(SupportedContracts[pauseContract.Hash].AbiHash) == 0 {
		err := fmt.Errorf("doPause: contarct with hash %s is not supported", pauseContract.Hash)
		log.Error(err)
		return nil, err
	}

	txOutMsg, err := h.txOutputProducer.PauseContract(ctx, pauseContract.Chain, pauseContract.Hash)
	if err != nil {
		log.Error("cannot get txOut for pausing contract, err = ", err)
		return nil, nil
	}

	h.signTx(ctx, txOutMsg.Data)

	return nil, nil
}

// signTx sends a TxOut to dheart for TSS signing.
func (h *HandlerPause) signTx(ctx sdk.Context, tx *types.TxOut) {
	log.Info("Delivering TXOUT for chain", tx.OutChain, " tx hash = ", tx.OutHash)
	if tx.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		log.Info("This TxOut is a contract deployment")
	}

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

	// 4. Send it to Dheart for signing.
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

func (h *HandlerPause) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
