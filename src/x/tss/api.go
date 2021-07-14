package tss

type Api struct {
}

func NewApi() *Api {
	return &Api{}
}

func (a *Api) Version() string {
	return "1.0"
}

func (a *Api) KeygenResult(chain string) bool {
	return true
}
