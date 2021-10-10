package types

import (
	"github.com/sisu-network/cosmos-sdk/codec"
	cdctypes "github.com/sisu-network/cosmos-sdk/codec/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"

	"github.com/sisu-network/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &EthTx{})
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
