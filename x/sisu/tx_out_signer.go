package sisu

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	etypes "github.com/ethereum/go-ethereum/core/types"
	hTypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type txOutSigner struct {
	keeper       keeper.Keeper
	partyManager PartyManager
	dheartClient tssclients.DheartClient
}

func NewTxOutSigner(keeper keeper.Keeper, partyManager PartyManager,
	dheartClient tssclients.DheartClient) *txOutSigner {
	return &txOutSigner{
		keeper:       keeper,
		partyManager: partyManager,
		dheartClient: dheartClient,
	}
}

func (q *txOutSigner) signTxOut(ctx sdk.Context, txOut *types.TxOut) {
	if libchain.IsETHBasedChain(txOut.OutChain) {
		q.signEthTx(ctx, txOut)
	}

	if libchain.IsCardanoChain(txOut.OutChain) {
		q.signCardanoTx(ctx, txOut)
	}
}

// signEthTx sends a TxOut to dheart for TSS signing.
func (q *txOutSigner) signEthTx(ctx sdk.Context, tx *types.TxOut) {
	log.Info("Delivering TXOUT for chain ", tx.OutChain, " tx hash = ", tx.OutHash)
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
				Id:          q.getKeysignRequestId(tx.OutChain, ctx.BlockHeight(), tx.OutHash),
				OutChain:    tx.OutChain,
				OutHash:     tx.OutHash,
				BytesToSign: hash[:],
			},
		},
	}

	pubKeys := q.partyManager.GetActivePartyPubkeys()

	err := q.dheartClient.KeySign(keysignReq, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func (q *txOutSigner) signCardanoTx(ctx sdk.Context, txOut *types.TxOut) {
	tx := &cardano.Tx{}
	if err := tx.UnmarshalCBOR(txOut.OutBytes); err != nil {
		log.Error("error when unmarshalling cardano tx out: ", err)
		return
	}

	txHash, err := tx.Hash()
	if err != nil {
		log.Error("error when getting cardano tx hash: ", err)
		return
	}

	signRequest := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_EDDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          q.getKeysignRequestId(txOut.OutChain, ctx.BlockHeight(), txOut.OutHash),
				OutChain:    txOut.OutChain,
				OutHash:     txOut.OutHash,
				BytesToSign: txHash[:],
			},
		},
	}

	pubKeys := q.partyManager.GetActivePartyPubkeys()
	err = q.dheartClient.KeySign(signRequest, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func (q *txOutSigner) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
