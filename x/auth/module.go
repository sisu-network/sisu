package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/gorilla/mux"
	abci "github.com/sisu-network/tendermint/abci/types"
	"github.com/spf13/cobra"

	"github.com/sisu-network/cosmos-sdk/client"
	"github.com/sisu-network/cosmos-sdk/codec"
	codectypes "github.com/sisu-network/cosmos-sdk/codec/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/cosmos-sdk/types/module"
	simtypes "github.com/sisu-network/cosmos-sdk/types/simulation"
	"github.com/sisu-network/cosmos-sdk/x/auth/client/cli"
	"github.com/sisu-network/cosmos-sdk/x/auth/client/rest"
	"github.com/sisu-network/cosmos-sdk/x/auth/keeper"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the auth module.
type AppModuleBasic struct{}

// Name returns the auth module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterLegacyAminoCodec registers the auth module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the auth
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the auth module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}

	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the auth module.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(clientCtx, rtr, StoreKey)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the auth module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	RegisterQueryHandlerClient(context.Background(), mux, NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the auth module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the root query command for the auth module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces registers interfaces and implementations of the auth module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

//____________________________________________________________________________

// AppModule implements an application module for the auth module.
type AppModule struct {
	AppModuleBasic

	cosmosAppModule CosmosAppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Marshaler, accountKeeper keeper.AccountKeeper, randGenAccountsFn RandomGenesisAccountsFn) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},

		cosmosAppModule: NewCosmosAppModule(cdc, accountKeeper, randGenAccountsFn),
	}
}

// Name returns the auth module's name.
func (app AppModule) Name() string {
	return app.cosmosAppModule.Name()
}

// RegisterInvariants performs a no-op.
func (app AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	app.cosmosAppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the auth module.
func (app AppModule) Route() sdk.Route {
	return app.cosmosAppModule.Route()
}

// QuerierRoute returns the auth module's querier route name.
func (app AppModule) QuerierRoute() string {
	return app.cosmosAppModule.QuerierRoute()
}

// LegacyQuerierHandler returns the auth module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return am.cosmosAppModule.LegacyQuerierHandler(legacyQuerierCdc)
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	am.cosmosAppModule.RegisterServices(cfg)
}

// InitGenesis performs genesis initialization for the auth module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, cdc, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the auth
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	return am.ExportGenesis(ctx, cdc)
}

// BeginBlock returns the begin blocker for the auth module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the auth module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the auth module
func (am AppModule) GenerateGenesisState(simState *module.SimulationState) {
	am.cosmosAppModule.GenerateGenesisState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized auth param changes for the simulator.
func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return am.cosmosAppModule.RandomizedParams(r)
}

// RegisterStoreDecoder registers a decoder for auth module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	am.cosmosAppModule.RegisterStoreDecoder(sdr)
}

// WeightedOperations doesn't return any auth module operation.
func (AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
