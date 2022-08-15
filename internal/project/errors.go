package project

type ErrUnknown struct {
	Err error
}

func (e ErrUnknown) Error() string {
	return e.Err.Error()
}

func (e ErrUnknown) Unwrap() error {
	return e.Err
}

type ErrProjectNotFound struct {
	Err error
}

func (e ErrProjectNotFound) Error() string {
	return e.Err.Error()
}

func (e ErrProjectNotFound) Unwrap() error {
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
