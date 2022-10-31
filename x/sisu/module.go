package sisu

import (
	"encoding/json"
	"fmt"
	"strings"

	// this line is used by starport scaffolding # 1

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/client/cli"
	"github.com/sisu-network/sisu/x/sisu/client/rest"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the capability module.
type AppModuleBasic struct {
	cdc codec.Marshaler
}

func NewAppModuleBasic(cdc codec.Marshaler) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the capability module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns the capability module's default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the capability module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterRESTRoutes registers the capability module's REST service handlers.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// this line is used by starport scaffolding # 2
}

// GetTxCmd returns the capability module's root tx command.
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the capability module's root query command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the capability module.
type AppModule struct {
	AppModuleBasic

	sisuHandler     *SisuHandler
	externalHandler *rest.ExternalHandler
	keeper          keeper.Keeper
	processor       *ApiHandler
	appKeys         common.AppKeys
	txSubmit        common.TxSubmit
	globalData      common.GlobalData
	valsManager     ValidatorManager
	txTracker       TxTracker
	txOutSigner     *txOutSigner
	mc              ManagerContainer
}

func NewAppModule(cdc codec.Marshaler,
	sisuHandler *SisuHandler,
	keeper keeper.Keeper,
	apiHandler *ApiHandler,
	valsManager ValidatorManager,
	mc ManagerContainer,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		sisuHandler:    sisuHandler,
		txSubmit:       mc.TxSubmit(),
		processor:      apiHandler,
		keeper:         keeper,
		appKeys:        mc.AppKeys(),
		globalData:     mc.GlobalData(),
		valsManager:    valsManager,
		txTracker:      mc.TxTracker(),
		txOutSigner:    NewTxOutSigner(mc.Keeper(), mc.PartyManager(), mc.DheartClient()),
		mc:             mc,
	}
}

// Name returns the capability module's name.
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// Route returns the capability module's message routing key.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, am.sisuHandler.NewHandler(am.processor, am.valsManager))
}

// QuerierRoute returns the capability module's query routing key.
func (AppModule) QuerierRoute() string { return types.QuerierRoute }

// LegacyQuerierHandler returns the capability module's Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterTssQueryServer(cfg.QueryServer(), keeper.NewGrpcQuerier(am.keeper))
}

// RegisterInvariants registers the capability module's invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the capability module's genesis initialization.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	// Saves initial token configs from genesis file.
	tokenIds := make([]string, 0)
	m := make(map[string]*types.Token)
	for _, token := range genState.Tokens {
		m[token.Id] = token
		tokenIds = append(tokenIds, token.Id)
	}
	am.keeper.SetTokens(ctx, m)
	log.Info("Tokens in the genesis file: ", strings.Join(tokenIds, ", "))

	// Save initial chain data
	chains := make([]string, 0)
	for _, chain := range genState.Chains {
		am.keeper.SaveChain(ctx, chain)
		chains = append(chains, chain.Id)
	}
	log.Info("Chains in the genesis file: ", strings.Join(chains, ", "))

	// Save vaults
	am.keeper.SetVaults(ctx, genState.Vaults)

	// Save params
	params := genState.Params
	am.keeper.SaveParams(ctx, params)
	savedParams := am.keeper.GetParams(ctx)
	log.Info("Tss params: ", savedParams)

	// Save Checkpoints
	for _, checkpoint := range genState.Checkpoints {
		am.keeper.AddGatewayCheckPoint(ctx, checkpoint)
	}

	// Create validator nodes
	valsMgr := am.valsManager
	validators := make([]abci.ValidatorUpdate, len(genState.Nodes))
	for i, node := range genState.Nodes {
		if !node.IsValidator {
			continue
		}

		pk, err := utils.GetCosmosPubKey(node.ConsensusKey.Type, node.ConsensusKey.Bytes)
		if err != nil {
			panic(err)
		}

		validators[i] = abci.Ed25519ValidatorUpdate(pk.Bytes(), 100)
		valsMgr.AddValidator(ctx, node)
	}

	// Reload data after reading the genesis
	am.mc.TransferQueue().Start(ctx)

	return validators
}

// ExportGenesis returns the capability module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// BeginBlock executes all ABCI BeginBlock logic respective to the capability module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	log.Verbose("BeginBlock, height = ", ctx.BlockHeight())

	if !am.globalData.AppInitialized() {
		cloneCtx := utils.CloneSdkContext(ctx)
		am.mc.TransferQueue().Start(cloneCtx)
		am.globalData.SetAppInitialized()
	}

	am.beginBlock(ctx, req.Header.Height)
}

// EndBlock executes all ABCI EndBlock logic respective to the capability module. It
// returns no validator updates.
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	log.Verbose("End block reached, height = ", ctx.BlockHeight())

	am.processor.EndBlock(ctx)

	am.txTracker.CheckExpiredTransaction()

	cloneCtx := utils.CloneSdkContext(ctx)
	am.globalData.SetReadOnlyContext(cloneCtx)

	params := am.keeper.GetParams(ctx)

	// Update gas price
	if ctx.BlockHeight()%UpdateGasPriceFrequency == 0 {
		// Get the gas price from deyes
		go func(params *types.Params) {
			gasPrices, err := am.processor.deyesClient.GetGasPrices(params.SupportedChains)
			if err != nil {
				log.Error("cannot get gas price from deyes")
			} else {
				msg := types.NewGasPriceMsg(am.appKeys.GetSignerAddress().String(),
					params.SupportedChains, ctx.BlockHeight(), gasPrices)
				am.txSubmit.SubmitMessageAsync(msg)
			}
		}(params)
	}

	// Process pending transfers
	am.mc.TransferQueue().ProcessTransfers(ctx)

	// Sign tx outs
	am.signTxOut(ctx)

	return []abci.ValidatorUpdate{}
}

func (am AppModule) signTxOut(ctx sdk.Context) {
	params := am.keeper.GetParams(ctx)
	height := ctx.BlockHeight()

	for _, chain := range params.SupportedChains {
		pendingInfo := am.keeper.GetPendingTxOutInfo(ctx, chain)
		if pendingInfo != nil {
			if pendingInfo.ExpiredBlock < height {
				log.Infof("Pending tx on chain %s expired. Clearing the pending tx.", chain)
				am.keeper.SetPendingTxOutInfo(ctx, chain, nil)

				// TODO: Put this back to the failure queue
				// queue := am.keeper.GetTxOutQueue(ctx, chain)
				// queue = append(queue, pendingInfo.TxOut)
				// am.keeper.SetTxOutQueue(ctx, chain, queue)
			}

			continue
		}

		queue := am.keeper.GetTxOutQueue(ctx, chain)
		if len(queue) == 0 {
			continue
		}

		txOut := queue[0]
		am.keeper.SetPendingTxOutInfo(ctx, txOut.Content.OutChain, &types.PendingTxOutInfo{
			TxOut: txOut,
			// ExpiredBlock: height + params.PendingTxTimeoutHeights[i],
			// TODO: Make this height configurable
			ExpiredBlock: height + 50,
		})
		am.keeper.SetTxOutQueue(ctx, txOut.Content.OutChain, queue[1:])

		if !am.globalData.IsCatchingUp() {
			log.Verbose("Signing txout hash = ", txOut.Content.OutHash)
			am.txOutSigner.signTxOut(ctx, txOut)
		}
	}
}
