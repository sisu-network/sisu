package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (sh *SisuHandler) NewHandler(processor *Processor, valsManager ValidatorManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {

		ctx = ctx.WithEventManager(sdk.NewEventManager())
		mc := sh.mc

		signers := msg.GetSigners()
		if len(signers) == 0 {
			return &sdk.Result{}, fmt.Errorf("empty signer for msg: %s", msg.String())
		}

		if RequireSubmitPerm(msg) && !valsManager.IsNodeAccount(ctx, signers[0]) {
			return &sdk.Result{}, fmt.Errorf("require permision: %s", msg.String())
		}

		switch msg := msg.(type) {
		case *types.KeygenWithSigner:
			return NewHandlerKeygen(mc).DeliverMsg(ctx, msg)
		case *types.KeygenResultWithSigner:
			return NewHandlerKeygenResult(mc).DeliverMsg(ctx, msg)
		case *types.ReshareResultWithSigner:
			return NewHandlerReshareResult(mc).DeliverMsg(ctx, msg)
		case *types.ContractsWithSigner:
			return NewHandlerContract(mc).DeliverMsg(ctx, msg)
		case *types.TxInWithSigner:
			return NewHandlerTxIn(mc).DeliverMsg(ctx, msg)
		case *types.TxOutWithSigner:
			return NewHandlerTxOut(mc).DeliverMsg(ctx, msg)
		case *types.TxOutContractConfirmWithSigner:
			return NewHandlerTxOutContractConfirmation(mc).DeliverMsg(ctx, msg)
		case *types.KeysignResult:
			return NewHandlerKeysignResult(mc).DeliverMsg(ctx, msg)
		case *types.GasPriceMsg:
			return NewHandlerGasPrice(mc).DeliverMsg(ctx, msg)
		case *types.UpdateTokenPrice:
			return NewHandlerTokenPrice(mc).DeliverMsg(ctx, msg)
		case *types.PauseContractMsg:
			return NewHandlerPauseContract(mc).DeliverMsg(ctx, msg)
		case *types.ResumeContractMsg:
			return NewHandlerResumeContract(mc).DeliverMsg(ctx, msg)
		case *types.ChangeOwnershipContractMsg:
			return NewHandlerContractChangeOwnership(mc).DeliverMsg(ctx, msg)
		case *types.ChangeLiquidPoolAddressMsg:
			return NewHandlerContractSetLiquidityAddress(mc).DeliverMsg(ctx, msg)
		case *types.LiquidityWithdrawFundMsg:
			return NewHandlerContractLiquidityWithdrawFund(mc).DeliverMsg(ctx, msg)
		case *types.ChangeValidatorSetMsg:
			return NewHandlerChangeValidatorSet(mc).DeliverMsg(ctx, msg)
		case *types.DepositSisuTokenMsg:
			return NewHandlerDepositSisuToken(mc).DeliverMsg(ctx, msg)
		case *types.SlashValidatorMsg:
			return NewHandlerSlashValidator(mc).DeliverMsg(ctx, msg)
		case *types.SetDheartIpAddressMsg:
			return NewHandlerSetDheartIPAddress(mc).DeliverMsg(ctx, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return &sdk.Result{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func RequireSubmitPerm(msg sdk.Msg) bool {
	switch msg.(type) {
	case *types.DepositSisuTokenMsg, *types.SetDheartIpAddressMsg, *types.SlashValidatorMsg:
		return false
	default:
		return true
	}
}
