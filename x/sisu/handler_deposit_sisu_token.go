package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

var BondAddr = sdk.AccAddress("0xbondaddress")

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
		return nil, nil
	}

	if err := h.doDepositSisuToken(ctx, msg); err != nil {
		return &sdk.Result{}, err
	}

	h.keeper.ProcessTxRecord(ctx, hash)
	return nil, nil
}

func (h *HandlerDepositSisuToken) doDepositSisuToken(ctx sdk.Context, msg *types.DepositSisuTokenMsg) error {
	depositAmt := msg.Data.Amount
	if err := h.mc.BankKeeper().SendCoinsFromAccountToModule(ctx, msg.GetSender(), BondName, sdk.Coins{
		sdk.NewCoin(common.SisuCoinName, sdk.NewInt(depositAmt)),
	}); err != nil {
		log.Error("error when send coin to bond addr. error = ", err)
		return err
	}

	return h.keeper.IncBalance(ctx, msg.GetSender(), depositAmt)
}
