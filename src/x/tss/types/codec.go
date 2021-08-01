package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenProposal{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenProposalVote{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenResult{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &ObservedTxs{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TxOut{})
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
