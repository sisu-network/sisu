package types

type KeygenEntity struct {
	Chain   string `json:"chain,omitempty"`
	Address string `json:"address,omitempty"`
	PubKey  []byte `json:"pub_key,omitempty"`
	Status  string `json:"status,omitempty"`
}
