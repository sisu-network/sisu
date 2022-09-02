package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenWithSigner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenResultWithSigner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TransferOutsMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TxOutMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TxOutResultMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &GasPriceMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &UpdateTokenPrice{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeysignResult{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &PauseContractMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &ResumeContractMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &ChangeOwnershipContractMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &LiquidityWithdrawFundMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &BlockHeightMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TransferBatchMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TransferFailureMsg{})

	registry.RegisterInterface("cosmos.crypto.PubKey", (*ctypes.PubKey)(nil))
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &ed25519.PubKey{})
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &multisig.LegacyAminoPubKey{})
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
