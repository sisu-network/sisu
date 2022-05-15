package sisu

import (
	"encoding/base64"
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
	marshaller := simapp.MakeTestEncodingConfig().Marshaler

	maccPerms := make(map[string][]string)
	maccPerms[BondName] = []string{}
	maccPerms[authtypes.Burner] = []string{authtypes.Burner}
	maccPerms[authtypes.Minter] = []string{authtypes.Minter}

	pk := paramskeeper.NewKeeper(marshaller, legacyCodec, testStoreKeys[paramstypes.StoreKey], sdk.NewTransientStoreKey(paramstypes.TStoreKey))
	ak := authkeeper.NewAccountKeeper(marshaller, testStoreKeys[authtypes.StoreKey], pk.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms)
	bk := bankkeeper.NewBaseKeeper(marshaller, testStoreKeys[banktypes.StoreKey], ak, pk.Subspace(banktypes.ModuleName), nil)

	return ak, bk
}

func TestHandlerDepositSisuToken_DepositToken(t *testing.T) {
	t.Parallel()

	initialPower := int64(100)
	initTokens := sdk.TokensFromConsensusPower(initialPower)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))

	t.Run("User can deposit", func(t *testing.T) {
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

		accBalance := int64(100)
		deposit := int64(10)

		account := appKeys.GetSignerAddress()
		bankKeeper.AddCoins(ctx, account, []sdk.Coin{sdk.NewCoin(common.SisuCoinName, sdk.NewInt(accBalance))})

		mc := MockManagerContainer(validatorManager, keeper, bankKeeper)
		handler := NewHandlerDepositSisuToken(mc)

		depositMsg := types.NewDepositSisuTokenMsg(appKeys.GetSignerAddress().String(), base64.StdEncoding.EncodeToString([]byte("pubkey")), deposit, 0)

		_, err := handler.DeliverMsg(ctx, depositMsg)
		require.Nil(t, err)

		balance, err := keeper.GetBondBalance(ctx, appKeys.GetSignerAddress())
		require.Nil(t, err)

		require.Equal(t, deposit, balance)
	})

	t.Run("User cannot deposit when balance is not enough", func(t *testing.T) {
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

		accBalance := int64(100)
		deposit := int64(200)

		account := appKeys.GetSignerAddress()
		bankKeeper.AddCoins(ctx, account, []sdk.Coin{sdk.NewCoin(common.SisuCoinName, sdk.NewInt(accBalance))})

		mc := MockManagerContainer(validatorManager, keeper, bankKeeper)
		handler := NewHandlerDepositSisuToken(mc)

		depositMsg := types.NewDepositSisuTokenMsg(appKeys.GetSignerAddress().String(), base64.StdEncoding.EncodeToString([]byte("pubkey")), deposit, 0)

		_, err := handler.DeliverMsg(ctx, depositMsg)
		require.NotNil(t, err)
	})
}
