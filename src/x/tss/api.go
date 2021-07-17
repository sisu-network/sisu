package tss

import (
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
	tTypes "github.com/sisu-network/tuktuk/types"
)

type Api struct {
	processor *Processor
	keeper    *keeper.Keeper
}

func NewApi(processor *Processor, keeper *keeper.Keeper) *Api {
	return &Api{
		processor: processor,
	}
}

func (a *Api) Version() string {
	return "1.0"
}

func (a *Api) KeygenResult(result tTypes.KeygenResult) bool {
	utils.LogInfo("There is a TSS Result")

	a.processor.OnKeygenResult(result)

	return true
}
