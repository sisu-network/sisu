package sisu

import (
	"fmt"

	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/components"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type SisuHandler struct {
	mc background.ManagerContainer
}

func NewSisuHandler(mc background.ManagerContainer) *SisuHandler {
	return &SisuHandler{
		mc: mc,
	}
}

func (sh *SisuHandler) NewHandler(processor *ApiHandler, valsManager components.ValidatorManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			log.Error("Signers length must be 1. Actual length = ", len(signers))
			return nil, fmt.Errorf("incorrect signers length: %d", len(signers))
		}

		if !valsManager.IsValidator(ctx, signers[0].String()) {
			log.Verbose("sender is not a validator ", signers[0].String())
			return nil, fmt.Errorf("sender is not a validator: %s", signers[0].String())
		}

		ctx = ctx.WithEventManager(sdk.NewEventManager())
		mc := sh.mc

		switch msg := msg.(type) {
		case *types.KeygenWithSigner:
			return NewHandlerKeygen(mc).DeliverMsg(ctx, msg)

		case *types.KeygenResultWithSigner:
			return NewHandlerKeygenResult(mc).DeliverMsg(ctx, msg)

		case *types.TxOutMsg:
			return NewHandlerTxOutProposal(mc).DeliverMsg(ctx, msg)

		case *types.TxOutResultMsg:
			return NewHandlerTxOutResult(mc).DeliverMsg(ctx, msg)

		case *types.KeysignResultMsg:
			return NewHandlerKeysignResult(mc).DeliverMsg(ctx, msg)

		case *types.BlockHeightMsg:
			return NewHandlerBlockHeight(mc.Keeper()).DeliverMsg(ctx, msg)

		case *types.TransferFailureMsg:
			return NewHanlderTransferFailure(mc.Keeper(), mc.PostedMessageManager()).DeliverMsg(ctx, msg)

		case *types.TxInMsg:
			return NewHandlerTxIn(mc.PostedMessageManager(), mc.Keeper(),
				mc.GlobalData(), mc.BridgeManager(), mc.ValidatorManager(), mc.PrivateDb()).DeliverMsg(ctx, msg)

		case *types.TransferRetryMsg:
			return NewHandlerTransferRetry(mc.PostedMessageManager(), mc.Keeper(),
				mc.GlobalData(), mc.BridgeManager(), mc.ValidatorManager(), mc.PrivateDb()).DeliverMsg(ctx, msg)

		case *types.TxOutVoteMsg:
			return NewHandlerTxOutConsensed(mc.PostedMessageManager(), mc.Keeper(),
				mc.PrivateDb()).DeliverMsg(ctx, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
