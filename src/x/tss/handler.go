package tss

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper, txSubmit common.TxSubmit, processor *Processor) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.KeygenProposal:
			return handleKeygenProposal(msg, processor)
		case *types.KeygenProposalVote:
			return handleKeygenProposalVote(msg, processor)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleKeygenProposal(msg *types.KeygenProposal, processor *Processor) (*sdk.Result, error) {
	data, err := processor.DeliverKeyGenProposal(msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleKeygenProposalVote(msg *types.KeygenProposalVote, processor *Processor) (*sdk.Result, error) {
	utils.LogDebug("Handling keygen proposal vote....")

	data, err := processor.DeliverKeyGenProposalVote(msg)
	return &sdk.Result{
		Data: data,
	}, err
}
