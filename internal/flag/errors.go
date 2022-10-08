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
	Message string
}

func (e ErrFlagNotFound) Error() string {
	return e.Message
}

type ErrInvalidData struct {
	Message string
}

func (e ErrInvalidData) Error() string {
	return e.Message
}

type ErrKeyNotUnique struct {
	Message string
}

func (e ErrKeyNotUnique) Error() string {
	return e.Message
}
