package flag

type ErrUnknown struct {
	Err error
}

func (e ErrUnknown) Error() string {
	return e.Err.Error()
}

func (e ErrUnknown) Unwrap() error {
	return e.Err
}

type ErrFlagNotFound struct {
	Err error
}

func (e ErrFlagNotFound) Error() string {
	return e.Err.Error()
}

func (e ErrFlagNotFound) Unwrap() error {
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

type ErrKeyNotUnique struct {
	Err error
}

func (e ErrKeyNotUnique) Error() string {
	return e.Err.Error()
}

func (e ErrKeyNotUnique) Unwrap() error {
	return e.Err
}
