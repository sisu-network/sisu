package keeper

import (
	"github.com/sisu-network/sisu/config"
)

type Keeper struct {
	tssConfig *config.TssConfig
}

func NewKeeper(tssConfig *config.TssConfig) *Keeper {
	return &Keeper{
		tssConfig: tssConfig,
	}
}
