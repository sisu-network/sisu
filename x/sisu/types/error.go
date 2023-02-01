package types

type ErrNotImplemented struct {
	msg string
}

func NewErrNotImplemented(msg string) *ErrNotImplemented {
	return &ErrNotImplemented{
		msg: msg,
	}
}

func (e *ErrNotImplemented) Error() string {
	if len(e.msg) > 0 {
		return e.msg
	}

	return "Method not implemented"
}
