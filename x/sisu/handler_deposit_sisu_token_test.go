package sisu

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
)

func testCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	banktypes.RegisterLegacyAminoCodec(cdc)
	authtypes.RegisterLegacyAminoCodec(cdc)
	types.RegisterCodec(cdc)
	return cdc
}

func getTestAccountAndBankKeepers() (authkeeper.AccountKeeper, bankkeeper.Keeper) {
	legacyCodec := testCodec()
	marshaler := simapp.MakeTestEncodingConfig().Marshaler

	maccPerms := make(map[string][]string)
	maccPerms[BondName] = []string{}
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}

	pk := paramskeeper.NewKeeper(marshaler, legacyCodec, testStoreKeys[paramstypes.StoreKey], sdk.NewTransientStoreKey(paramstypes.TStoreKey))
	ak := authkeeper.NewAccountKeeper(marshaler, testStoreKeys[authtypes.StoreKey], pk.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms)
	bk := bankkeeper.NewBaseKeeper(marshaler, testStoreKeys[banktypes.StoreKey], ak, pk.Subspace(banktypes.ModuleName), nil)

	return ak, bk
}

func TestHandlerDepositSisuToken_DepositToken(t *testing.T) {
	t.Parallel()

	initialPower := int64(100)
	initTokens := sdk.TokensFromConsensusPower(initialPower)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))

	t.Run("emtpy", func(t *testing.T) {
		t.Parallel()

		ctx := testContext()
		authKeeper, bankKeeper := getTestAccountAndBankKeepers()
		keeper := keeperTestGenesis(ctx)
		validatorManager := NewValidatorManager(keeper)
		appKeys := common.NewMockAppKeys()

		baseAcc := authKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("baseAcc"))
		bankKeeper.SetSupply(ctx, banktypes.NewSupply(initCoins))
		authKeeper.SetModuleAccount(ctx, authtypes.NewEmptyModuleAccount(BondName))
		authKeeper.SetAccount(ctx, baseAcc)

		account := appKeys.GetSignerAddress()
		bankKeeper.AddCoins(ctx, account, []sdk.Coin{sdk.NewCoin(common.SisuCoinName, sdk.NewInt(100))})

		mc := MockManagerContainer(validatorManager, keeper, bankKeeper)
		handler := NewHandlerDepositSisuToken(mc)

		depositMsg := types.NewDepositSisuTokenMsg(appKeys.GetSignerAddress().String(), "pubKey", 10, 0)

		_, err := handler.DeliverMsg(ctx, depositMsg)
		fmt.Println("err = ", err)
		require.NotNil(t, err)
	})
}
