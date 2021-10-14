package types

const (
	ContractStateDeploying = "deploying"
	ContractStateDeployed  = "deployed"
)

// Data model of contract which could be used for datbase.
type ContractEntity struct {
	Chain    string
	Hash     string
	ByteCode []byte
	Name     string
	Address  string
	Status   string
}
