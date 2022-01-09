package tss

import (
	"github.com/sisu-network/sisu/x/tss/keeper"
)

type ApiHandler struct {
	processor *Processor
	keeper    *keeper.DefaultKeeper
}

func NewApi(processor *Processor, keeper *keeper.DefaultKeeper) *ApiHandler {
	return &ApiHandler{
		processor: processor,
	}
}

func (a *ApiHandler) Version() string {
	return "1.0"
}

// Empty function for checking health only.
func (a *ApiHandler) CheckHealth() {
}
