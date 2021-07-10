package tss

type Api struct {
}

func NewApi() *Api {
	return &Api{}
}

func (a *Api) Version() string {
	return "1.0"
}
