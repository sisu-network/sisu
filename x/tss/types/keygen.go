package types

const (
	KEYGEN_STATUS_GENERATING = "generating"
	KEYGEN_STATUS_GENERATED  = "generated"
)

type KeygenEntity struct {
	Type       string
	Address    string
	Pubkey     []byte
	Status     string
	StartBlock int64
}
