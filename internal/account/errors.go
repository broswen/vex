package account

type ErrUnknown struct {
	Err error
}

func (e ErrUnknown) Error() string {
	return e.Err.Error()
}

func (e ErrUnknown) Unwrap() error {
	return e.Err
}

type ErrAccountNotFound struct {
	Err error
}

func (e ErrAccountNotFound) Error() string {
	return e.Err.Error()
}

func (e ErrAccountNotFound) Unwrap() error {
	return e.Err
}

type ErrInvalidData struct {
	Err error
}

func (e ErrInvalidData) Error() string {
	return e.Err.Error()
}

func (e ErrInvalidData) Unwrap() error {
	return e.Err
}
