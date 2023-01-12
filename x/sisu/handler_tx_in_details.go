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
	// Get the has from the thin TxIn.
	txInHash, _, err := keeper.GetTxRecordHash(&types.TxInMsg{
		Signer: msg.Signer,
		Data:   msg.Data.TxIn,
	})
	if err != nil {
		log.Errorf("Cannot get tx record hash, err = ", err)
		return &sdk.Result{}, nil
	}

	fmt.Println("TxIn Hash 2 = ", hex.EncodeToString(txInHash))
	shouldProcess, _ := h.pmm.ShouldProcessMsg(ctx, msg)

	fmt.Println("AAAAA Inside handler tx in details")

	assignedNode := h.valsManager.GetAssignedValidator(ctx, msg.Data.TxIn.Id)
	// Case 1: the thin tx is confirmed but no tx details is saved yet.
	if h.keeper.IsTxRecordProcessed(ctx, txInHash) {
		fmt.Println("AAAAA 0000")
		if shouldProcess || assignedNode.AccAddress == msg.Signer {
			// TODO: Do verification for this message before saving it (verify that all transfer data
			// is correct and match the TxIn transaction)
			doTxInDetails(ctx, h.keeper, msg)
		}
	} else {
		fmt.Println("AAAAA 111111")
		// Case 2: thin tx is not confirmed yet. We only want to save the TxInDetails of the assigned
		// node for this TxIn. Later on, when the thin TxIn is confirmed, we already have the details
		// for it to process.
		log.Verbosef("Assigned node = %s, msg signer = %s", assignedNode.AccAddress, msg.Signer)
		if assignedNode.AccAddress == msg.Signer {
			// TODO: Do verification for this message before saving it (verify that all transfer data
			// is correct and match the TxIn transaction)
			h.keeper.SetTxInDetails(ctx, msg.Data.TxIn.Id, msg.Data)
		}
	}

	return &sdk.Result{}, nil
}

// doTxInDetails should only be called when majority of nodes has submitted thin TxIn (to confirm)
// and either:
// 1) The assigned validator submitted the TxInDetails within a time frame
// 2) The assigned validator fails to submit TxInDetails but majority of nodes have submitted the
// TxIn details.
func doTxInDetails(ctx sdk.Context, k keeper.Keeper, msg *types.TxInDetailsMsg) {
	log.Verbosef("Process TxInDetails with TxIn id %s", msg.Data.TxIn.Id)

	hash, _, err := keeper.GetTxRecordHash(msg)
	if err != nil {
		log.Errorf("doTxInDetails: Failed to get tx record hash for TxInDetailsMsg")
		return
	}

	k.ProcessTxRecord(ctx, hash)

	// 1 .Save the tx in details.
	k.SetTxInDetails(ctx, msg.Data.FromChain, msg.Data)

	// 2. Save the transfers
	saveTransfers(ctx, k, msg.Data.Transfers)
}
