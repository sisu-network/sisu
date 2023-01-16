package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	bin "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
	etypes "github.com/sisu-network/deyes/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	dhtypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (a *ApiHandler) processETHSigningResult(ctx sdk.Context, result *dhtypes.KeysignResult,
	signMsg *dhtypes.KeysignMessage, index int) error {
	// Find the tx in txout table
	txOut := a.keeper.GetTxOut(ctx, signMsg.OutChain, signMsg.OutHash)
	if txOut == nil {
		return fmt.Errorf("Cannot find tx out with hash %s", signMsg.OutHash)
	}

	tx := &ethtypes.Transaction{}
	if err := tx.UnmarshalBinary(result.Request.KeysignMessages[index].Bytes); err != nil {
		return fmt.Errorf("cannot unmarshal tx, err = %v", err)
	}

	// Create full tx with signature.
	chainId := libchain.GetChainIntFromId(signMsg.OutChain)
	if len(result.Signatures[index]) != 65 {
		log.Error("Signature length is not 65 for chain: ", chainId)
	}
	signedTx, err := tx.WithSignature(ethtypes.NewLondonSigner(chainId), result.Signatures[index])
	if err != nil {
		return fmt.Errorf("cannot set signature for tx, err = %v", err)
	}

	bz, err := signedTx.MarshalBinary()
	if err != nil {
		return fmt.Errorf("cannot marshal signedTx")
	}

	// // TODO: Broadcast the keysign result that includes this TxOutSig.
	// // Save this to TxOutSig
	log.Verbosef("ETH keysign result chain = %s, hash (no sig) = %s, hash (signed) = %s",
		signMsg.OutChain, signMsg.OutHash, signedTx.Hash().String())
	a.privateDb.SaveTxOutSig(&types.TxOutSig{
		Chain:       signMsg.OutChain,
		HashWithSig: signedTx.Hash().String(),
		HashNoSig:   signMsg.OutHash,
	})

	err = a.deploySignedTx(ctx, bz, signMsg.OutChain, signedTx.Hash().String())
	if err != nil {
		return fmt.Errorf("deployment error: %v", err)
	}

	return nil
}

func (a *ApiHandler) processCardanoSigningResult(ctx sdk.Context, result *dhtypes.KeysignResult,
	signMsg *dhtypes.KeysignMessage, index int) error {
	log.Info("Processing Cardano signing result ...")
	txOut := a.keeper.GetTxOut(ctx, signMsg.OutChain, signMsg.OutHash)
	if txOut == nil {
		err := fmt.Errorf("cannot find tx out with hash %s", signMsg.OutHash)
		log.Error(err)
		return err
	}

	tx := &cardano.Tx{}
	if err := tx.UnmarshalCBOR(txOut.Content.OutBytes); err != nil {
		log.Error("error when unmarshalling cardano tx: ", err)
		return err
	}

	pubkey := a.keeper.GetKeygenPubkey(ctx, libchain.GetKeyTypeForChain(signMsg.OutChain))
	if len(pubkey) == 0 {
		err := fmt.Errorf("cannot find pubkey for type %s", libchain.GetKeyTypeForChain(signMsg.OutChain))
		log.Error(err)
		return err
	}

	for i := range tx.WitnessSet.VKeyWitnessSet {
		tx.WitnessSet.VKeyWitnessSet[i] = cardano.VKeyWitness{VKey: pubkey, Signature: result.Signatures[index]}
	}

	hashWSig, err := tx.Hash()
	if err != nil {
		log.Error(err)
		return err
	}

	a.privateDb.SaveTxOutSig(&types.TxOutSig{
		Chain:       signMsg.OutChain,
		HashWithSig: hashWSig.String(),
		HashNoSig:   signMsg.OutHash,
	})

	txBytes, err := tx.MarshalCBOR()
	if err != nil {
		log.Error("error when marshal cardano tx: ", err)
		return err
	}
	hash, err := tx.Hash()
	if err != nil {
		return nil
	}

	err = a.deploySignedTx(ctx, txBytes, signMsg.OutChain, result.Request.KeysignMessages[index].OutHash)
	if err != nil {
		log.Error("deployment error: ", err)
		return err
	}

	log.Info("Sent signed cardano tx to deyes, tx hash = ", hash)

	return nil
}

func (a *ApiHandler) processSolanaKeysignResult(ctx sdk.Context, result *dhtypes.KeysignResult,
	signMsg *dhtypes.KeysignMessage, index int) error {
	tx := solanago.Transaction{}
	message := solanago.Message{}
	err := message.UnmarshalLegacy(bin.NewCompactU16Decoder(signMsg.BytesToSign))
	if err != nil {
		return err
	}

	// TODO: Support multi transactions here.
	tx.Message = message
	tx.Signatures = []solanago.Signature{solanago.SignatureFromBytes(result.Signatures[0])}

	txBytes, err := tx.MarshalBinary()
	if err != nil {
		return err
	}

	a.privateDb.SaveTxOutSig(&types.TxOutSig{
		Chain:       signMsg.OutChain,
		HashWithSig: tx.Signatures[0].String(),
		HashNoSig:   signMsg.OutHash,
	})

	log.Verbose("Sending signed solana tx to deyes....")
	err = a.deploySignedTx(ctx, txBytes, signMsg.OutChain, tx.Signatures[0].String())
	if err != nil {
		log.Error("deployment error: ", err)
		return err
	}

	return nil
}

// deploySignedTx creates a deployment request and sends it to deyes.
func (a *ApiHandler) deploySignedTx(ctx sdk.Context, bz []byte, outChain string, trackHash string) error {
	log.Verbose("Sending final tx to the deyes for deployment for chain ", outChain)

	pubkey := a.keeper.GetKeygenPubkey(ctx, libchain.GetKeyTypeForChain(outChain))
	if pubkey == nil {
		return fmt.Errorf("Cannot get pubkey for chain %s", outChain)
	}

	request := &etypes.DispatchedTxRequest{
		Chain:  outChain,
		TxHash: trackHash,
		Tx:     bz,
		PubKey: pubkey,
	}

	go func(request *eyesTypes.DispatchedTxRequest) {
		result, err := a.deyesClient.Dispatch(request)

		// Handle failure case.
		if err != nil || (result != nil && !result.Success) {
			log.Error("Deployment failed!, err = ", err)

			txOut := a.getTxOutFromSignedHash(outChain, trackHash)

			if txOut == nil {
				log.Errorf("Cannot find txOut for dispath result with signed hash = %s, chain = %s", trackHash, outChain)
				return
			}

			txOutId := txOut.GetId()
			// Report this as failure. Submit to the Sisu chain
			txOutResult := &types.TxOutResult{
				TxOutId:  txOutId,
				OutChain: txOut.Content.OutChain,
				OutHash:  txOut.Content.OutHash,
			}
			txOutResult.Result = types.TxOutResultType_GENERIC_ERROR

			if result != nil {
				log.Verbose("Result error = ", result.Err)
				switch result.Err {
				case etypes.ErrNotEnoughBalance:
					txOutResult.Result = types.TxOutResultType_NOT_ENOUGH_NATIVE_BALANCE
				}
			}

			a.submitTxOutResult(txOutResult)
		} else {
			log.Verbose("Tx is sent to deyes!")
		}
	}(request)

	return nil
}
