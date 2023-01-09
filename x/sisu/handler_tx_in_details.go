package sisu

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxInDetails struct {
	pmm           PostedMessageManager
	keeper        keeper.Keeper
	globalData    common.GlobalData
	bridgeManager chains.BridgeManager
	valsManager   ValidatorManager
}

func NewHandlerTxInDetails(
	pmm PostedMessageManager,
	keeper keeper.Keeper,
	globalData common.GlobalData,
	bridgeManager chains.BridgeManager,
	valsManager ValidatorManager,
) *HandlerTxInDetails {
	return &HandlerTxInDetails{
		pmm:           pmm,
		keeper:        keeper,
		globalData:    globalData,
		bridgeManager: bridgeManager,
		valsManager:   valsManager,
	}
}

func (h *HandlerTxInDetails) DeliverMsg(ctx sdk.Context, msg *types.TxInDetailsMsg) (*sdk.Result, error) {
	txInHash, _, err := keeper.GetTxRecordHash(&types.TxInMsg{
		Signer: msg.Signer,
		Data:   msg.Data.TxIn,
	})
	if err != nil {
		log.Errorf("Cannot get tx record hash, err = ", err)
		return &sdk.Result{}, nil
	}

	fmt.Println("TxIn Hash 2 = ", hex.EncodeToString(txInHash))
	shouldProcess, hash := h.pmm.ShouldProcessMsg(ctx, msg)

	fmt.Println("AAAAA Inside handler tx in details")

	if h.keeper.IsTxRecordProcessed(ctx, txInHash) {
		fmt.Println("AAAAA 0000")

		// Case 1: the thin tx is confirmed but no tx details is saved yet.
		if shouldProcess {
			h.keeper.ProcessTxRecord(ctx, hash)

			// 1 .Save the tx in details.
			h.keeper.SetTxInDetails(ctx, msg.Data.FromChain, msg.Data)

			// 2. Save the transfers
			saveTransfers(ctx, h.keeper, msg.Data.Transfers)
		} else {
			// TODO: Check that this node is the assigned validator for the transaction and the
		}
	} else {
		fmt.Println("AAAAA 111111")
		// Case 2: thin tx is not confirmed yet. We only want to save the TxInDetails of the assigned
		// node for this TxIn. Later on, when the thin TxIn is confirmed, we already have the details
		// for it to process.
		assignedNode := h.valsManager.GetAssignedValidator(ctx, msg.Data.TxIn.Id)
		if assignedNode.AccAddress == msg.Signer {
			// TODO: Do verification for this message before saving it (verify that all transfer data
			// is correct and match the TxIn transaction)
			h.keeper.SetTxInDetails(ctx, msg.Data.FromChain, msg.Data)
		}
	}

	return &sdk.Result{}, nil
}
