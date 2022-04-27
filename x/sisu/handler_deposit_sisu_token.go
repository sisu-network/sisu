package sisu

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerDepositSisuToken struct {
	pmm    PostedMessageManager
	mc     ManagerContainer
	keeper keeper.Keeper
}

func NewHandlerDepositSisuToken(mc ManagerContainer) *HandlerDepositSisuToken {
	return &HandlerDepositSisuToken{
		pmm:    mc.PostedMessageManager(),
		mc:     mc,
		keeper: mc.Keeper(),
	}
}

func (h *HandlerDepositSisuToken) DeliverMsg(ctx sdk.Context, msg *types.DepositSisuTokenMsg) (*sdk.Result, error) {
	process, hash := h.pmm.ShouldProcessMsg(ctx, msg)
	if !process {
		return &sdk.Result{}, nil
	}

	if err := h.doDepositSisuToken(ctx, msg); err != nil {
		return &sdk.Result{}, err
	}

	h.keeper.ProcessTxRecord(ctx, hash)
	return &sdk.Result{}, nil
}

func (h *HandlerDepositSisuToken) doDepositSisuToken(ctx sdk.Context, msg *types.DepositSisuTokenMsg) error {
	b, err := h.keeper.GetBalance(ctx, msg.GetSender())
	if err != nil {
		return err
	}
	log.Debug("balance in keeper: ", b)
	depositAmt := msg.Data.Amount
	balance := h.mc.BankKeeper().GetBalance(ctx, msg.GetSender(), common.SisuCoinName)
	if balance.Amount.Int64() < depositAmt {
		err = errors.New(fmt.Sprintf("not enough sisu balance. Require %d, has %d", depositAmt, balance.Amount.Int64()))
		log.Error(err)
		return err
	}

	log.Debug("Balance before: ", balance)
	if err := h.mc.BankKeeper().SendCoinsFromAccountToModule(ctx, msg.GetSender(), BondName, sdk.Coins{
		sdk.NewCoin(common.SisuCoinName, sdk.NewInt(depositAmt)),
	}); err != nil {
		log.Error("error when send coin to bond addr. error = ", err)
		return err
	}

	return h.keeper.IncBalance(ctx, msg.GetSender(), depositAmt)
}
