package helper

import (
	"encoding/json"
	"os"

	"github.com/sisu-network/sisu/x/sisu/types"
)

func GetChains(file string) []*types.Chain {
	chains := []*types.Chain{}

	dat, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(dat, &chains); err != nil {
		panic(err)
	}

	return chains
}
