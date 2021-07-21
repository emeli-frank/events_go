package validation

type Error interface {
	error
	Info() string // todo:: fix, this was just placed here to differentiate this from the inbuilt error
}

// NewError returns validation.Error from a string
func NewError(str string) error {
	err := stringError(str)
	return &err
}

type stringError string

func (e *stringError) Error() string {
	return string(*e)
}

func (e *stringError) Info() string {
	return e.Error()
}
