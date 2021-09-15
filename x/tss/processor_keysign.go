package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dhTypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"

	etypes "github.com/sisu-network/dcore/core/types"
)

// This function is called after dheart sends Sisu keysign result.
func (p *Processor) OnKeysignResult(result *dhTypes.KeysignResult) {
	// Post the keysign result to cosmos chain.
	msg := types.NewKeysignResult(p.appKeys.GetSignerAddress().String(), result.OutChain, result.OutHash, result.Success, result.Signature)
	go p.txSubmit.SubmitMessage(msg)

	// Sends it to deyes for deployment.
	if result.Success {
		tx := &etypes.Transaction{}
		if err := tx.UnmarshalBinary(result.OutBytes); err != nil {
			utils.LogError("cannot unmarshal tx, err =", err)
			return
		}

		chainId := utils.GetChainIntFromId(result.OutChain)
		signedTx, err := tx.WithSignature(etypes.NewEIP2930Signer(chainId), result.Signature)
		if err != nil {
			utils.LogError("cannot set signatuer for tx, err =", err)
			return
		}

		utils.LogDebug("Sending final tx to the deyes for deployment")
		deyeClient := p.deyesClients[result.OutChain]
		if deyeClient != nil {
			bz, err := signedTx.MarshalBinary()
			if err != nil {
				utils.LogError("cannot marshal tx")
				return
			}

			deyeClient.Dispatch(result.OutChain, bz)
		} else {
			utils.LogError("Cannot find deyes client for chain", result.OutChain)
		}
	}
}

func (p *Processor) CheckKeysignResult(ctx sdk.Context, msg *types.KeysignResult) error {
	return nil
}

func (p *Processor) DeliverKeysignResult(ctx sdk.Context, msg *types.KeysignResult) ([]byte, error) {
	// TODO: implements this to handle blame.
	return nil, nil
}
