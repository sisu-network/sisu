package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type SisuHandler struct {
	mc ManagerContainer
}

func NewSisuHandler(mc ManagerContainer) *SisuHandler {
	return &SisuHandler{
		mc: mc,
	}
}

func (sh *SisuHandler) NewHandler(processor *ApiHandler, valsManager ValidatorManager) sdk.Handler {
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

		case *types.TransfersMsg:
			return NewHandlerTransfers(mc).DeliverMsg(ctx, msg)

		case *types.TxOutMsg:
			return NewHandlerTxOut(mc).DeliverMsg(ctx, msg)

		case *types.TxOutResultMsg:
			return NewHandlerTxOutResult(mc).DeliverMsg(ctx, msg)

		case *types.KeysignResult:
			return NewHandlerKeysignResult(mc).DeliverMsg(ctx, msg)

		case *types.GasPriceMsg:
			return NewHandlerGasPrice(mc).DeliverMsg(ctx, msg)

		case *types.UpdateTokenPrice:
			return NewHandlerTokenPrice(mc).DeliverMsg(ctx, msg)

		case *types.BlockHeightMsg:
			return NewHandlerBlockHeight(mc.Keeper()).DeliverMsg(ctx, msg)

		case *types.TransferFailureMsg:
			return NewHanlderTransferFailure(mc.Keeper(), mc.PostedMessageManager()).DeliverMsg(ctx, msg)

		case *types.UpdateSolanaRecentHashMsg:
			return NewHandlerUpdateSolanaRecentHash(mc.Keeper()).DeliverMsg(ctx, msg)

		case *types.AdjustEthNonceMsg:
			return NewHandlerAdjustEthNonce(
				mc.PostedMessageManager(), mc.Keeper(), mc.PrivateDb(),
			).DeliverMsg(ctx, msg)

		case *types.TxInMsg:
			return NewHandlerTxIn(mc.PostedMessageManager(), mc.Keeper(), mc.ValidatorManager(),
				mc.GlobalData()).DeliverMsg(ctx, msg)

		case *types.TxInDetailsMsg:
			return NewHandlerTxInDetails(mc.PostedMessageManager(), mc.Keeper(),
				mc.GlobalData(), mc.BridgeManager(), mc.ValidatorManager()).DeliverMsg(ctx, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
