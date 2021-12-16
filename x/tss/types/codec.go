package types

import (
	"github.com/sisu-network/cosmos-sdk/codec"
	cdctypes "github.com/sisu-network/cosmos-sdk/codec/types"
	"github.com/sisu-network/cosmos-sdk/crypto/keys/ed25519"
	"github.com/sisu-network/cosmos-sdk/crypto/keys/multisig"
	"github.com/sisu-network/cosmos-sdk/crypto/keys/secp256k1"
	ctypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"

	"github.com/sisu-network/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenProposalWithSigner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenResult{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &ObservedTx{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TxOut{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeysignResult{})

	registry.RegisterInterface("cosmos.crypto.PubKey", (*ctypes.PubKey)(nil))
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &ed25519.PubKey{})
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &multisig.LegacyAminoPubKey{})
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
