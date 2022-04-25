package sisu

import (
	"context"
	"errors"
	"fmt"
	"github.com/sisu-network/sisu/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/sisu-network/lib/log"
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
	resp, err := h.mc.BankKeeper().Balance(context.Background(), &bTypes.QueryBalanceRequest{
		Address: msg.Signer,
		Denom:   common.SisuCoinName,
	})

	if err != nil {
		log.Error("error when querying sisu balance. error = ", err)
		return err
	}

	balance := resp.Balance.Amount.Int64()
	depositAmt := msg.Data.Amount
	if balance < depositAmt {
		err = errors.New(fmt.Sprintf("not enough sisu token. Require %d, has %d", depositAmt, balance))
		log.Error(err)
		return err
	}

	if err = h.mc.BankKeeper().SendCoins(ctx, msg.GetSender(), BondAddr, sdk.Coins{
		sdk.NewCoin(common.SisuCoinName, sdk.NewInt(depositAmt)),
	}); err != nil {
		log.Error("error when send coin to bond addr. error = ", err)
		return err
	}

	return h.keeper.IncBalance(ctx, msg.GetSender(), depositAmt)
}
