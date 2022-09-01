package types

type MethodType int

const (
	MethodUnknown MethodType = iota
	MethodNativeTransfer
	MethodTransferOut
	MethodAddSpender
)
