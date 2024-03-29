package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenWithSigner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeygenResultWithSigner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TransfersMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TxOutMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TxOutVoteMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TxOutResultMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &KeysignResultMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &PauseContractMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &ResumeContractMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &ChangeOwnershipContractMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &LiquidityWithdrawFundMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &BlockHeightMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TransferFailureMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TxInMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &TransferRetryMsg{})

	registry.RegisterInterface("cosmos.crypto.PubKey", (*ctypes.PubKey)(nil))
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &ed25519.PubKey{})
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*ctypes.PubKey)(nil), &multisig.LegacyAminoPubKey{})
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
