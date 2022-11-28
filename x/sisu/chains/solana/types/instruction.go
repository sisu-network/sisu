package types

type InstructionType byte

const (
	Initialize    InstructionType = 0
	TransferOut                   = 1
	TranserIn                     = 2
	AddSpender                    = 3
	RemoveSpender                 = 4
	ChangeAdmin                   = 5
)
