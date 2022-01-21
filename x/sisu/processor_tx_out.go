package sisu

import (
	"fmt"
	"strconv"

	hTypes "github.com/sisu-network/dheart/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/types"

	etypes "github.com/ethereum/go-ethereum/core/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
)

func (p *Processor) deliverTxOut(ctx sdk.Context, signerMsg *types.TxOutWithSigner) ([]byte, error) {
	if process, hash := p.shouldProcessMsg(ctx, signerMsg); process {
		p.doTxOut(ctx, signerMsg)
		p.publicDb.ProcessTxRecord(hash)
	}

	return nil, nil
}

// deliverTxOut executes a TxOut transaction after it's included in Sisu block. If this node is
// catching up with the network, we would not send the tx to TSS for signing.
func (p *Processor) doTxOut(ctx sdk.Context, msgWithSigner *types.TxOutWithSigner) ([]byte, error) {
	txOut := msgWithSigner.Data

	log.Info("Delivering TxOut")

	// Save this to KVStore
	p.publicDb.SaveTxOut(txOut)

	// If this is a txout deployment,

	// Do key signing if this node is not catching up.
	if !p.globalData.IsCatchingUp() {
		// Only Deliver TxOut if the chain has been up to date.
		if libchain.IsETHBasedChain(txOut.OutChain) {
			p.signTx(ctx, txOut)
		}
	}

	return nil, nil
}

// signTx sends a TxOut to dheart for TSS signing.
func (p *Processor) signTx(ctx sdk.Context, tx *types.TxOut) {
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
				Id:          p.getKeysignRequestId(tx.OutChain, ctx.BlockHeight(), tx.OutHash),
				InChain:     tx.InChain,
				OutChain:    tx.OutChain,
				OutHash:     tx.OutHash,
				BytesToSign: hash[:],
			},
		},
	}

	pubKeys := p.partyManager.GetActivePartyPubkeys()

	err := p.dheartClient.KeySign(keysignReq, pubKeys)

	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func (p *Processor) getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
