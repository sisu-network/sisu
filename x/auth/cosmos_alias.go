// nolint
package auth

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	types "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	ModuleName                    = types.ModuleName
	StoreKey                      = types.StoreKey
	FeeCollectorName              = types.FeeCollectorName
	QuerierRoute                  = types.QuerierRoute
	DefaultMaxMemoCharacters      = types.DefaultMaxMemoCharacters
	DefaultTxSigLimit             = types.DefaultTxSigLimit
	DefaultTxSizeCostPerByte      = types.DefaultTxSizeCostPerByte
	DefaultSigVerifyCostED25519   = types.DefaultSigVerifyCostED25519
	DefaultSigVerifyCostSecp256k1 = types.DefaultSigVerifyCostSecp256k1
	QueryAccount                  = types.QueryAccount
)

var (
	// functions aliases
	NewCosmosAppModule          = auth.NewAppModule
	NewBaseAccount              = types.NewBaseAccount
	ProtoBaseAccount            = types.ProtoBaseAccount
	NewBaseAccountWithAddress   = types.NewBaseAccountWithAddress
	NewGenesisState             = types.NewGenesisState
	DefaultGenesisState         = types.DefaultGenesisState
	ValidateGenesis             = types.ValidateGenesis
	SanitizeGenesisAccounts     = types.SanitizeGenesisAccounts
	AddressStoreKey             = types.AddressStoreKey
	NewParams                   = types.NewParams
	ParamKeyTable               = types.ParamKeyTable
	DefaultParams               = types.DefaultParams
	ValidateGenAccounts         = types.ValidateGenAccounts
	GetGenesisStateFromAppState = types.GetGenesisStateFromAppState
	RegisterLegacyAminoCodec    = types.RegisterLegacyAminoCodec
	RegisterQueryHandlerClient  = types.RegisterQueryHandlerClient
	NewQueryClient              = types.NewQueryClient
	RegisterInterfaces          = types.RegisterInterfaces

	// variable aliases
	CosmosModuleCdc           = types.ModuleCdc
	AddressStoreKeyPrefix     = types.AddressStoreKeyPrefix
	GlobalAccountNumberKey    = types.GlobalAccountNumberKey
	KeyMaxMemoCharacters      = types.KeyMaxMemoCharacters
	KeyTxSigLimit             = types.KeyTxSigLimit
	KeyTxSizeCostPerByte      = types.KeyTxSizeCostPerByte
	KeySigVerifyCostED25519   = types.KeySigVerifyCostED25519
	KeySigVerifyCostSecp256k1 = types.KeySigVerifyCostSecp256k1
)

type (
	BaseAccount             = types.BaseAccount
	AccountRetriever        = types.AccountRetriever
	GenesisState            = types.GenesisState
	Params                  = types.Params
	GenesisAccountIterator  = types.GenesisAccountIterator
	RandomGenesisAccountsFn = types.RandomGenesisAccountsFn

	CosmosAppModuleBasic = auth.AppModuleBasic
	CosmosAppModule      = auth.AppModule
)
