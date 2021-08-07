package form

import (
	"events/pkg/validation"
	"net/url"
)









// Error is a collection of errors from validating form fields
type Error interface {
	Error() string
	ErrorMessages() map[string][]string
	Add(key, msg string)
}

type validationError map[string][]string // todo:: see if to convert to map[string]string

func (e *validationError) Error() string {
	return "Validation error" // todo:: print validation error messages serialized
}

func (e *validationError) ErrorMessages() map[string][]string {
	// because e can be nil, dereferencing it would panic so we do
	// a check first
	if e == nil {
		return map[string][]string{}
	}
	return *e
}

func (e *validationError) Add(key, msg string) {
	(*e)[key] = append((*e)[key], msg)
}

func (e *validationError) get(field string) string {
	msgs := (*e)[field]
	if len(msgs) < 1 {
		return ""
	}
	return msgs[0] // todo:: fix implementation
}













// New is a constructor function that returns a form instance
func New(data url.Values) *form {
	/*e :=validationError(map[string][]string{})
	return &form{
		data,
		&e,
		nil,
	}*/

	return &form{
		data,
		nil,
		nil,
	}
}

type form struct {
	url.Values
	ValidationErrors *validationError

	// Error is an unexpected error and not a validation error
	unexpectedErr error
}

/*func (f *form) Valid() bool {
	return len(f.ValidationErrors) == 0
}*/

func (f *form) Error() error {
	// this if statement might seem redundant but it is need to avoid unexpected behaviour
	// where f.ValidationErrors can be nil here but not nil in the main function.
	// todo:: learn why
	if f.ValidationErrors == nil {
		return nil
	}

	return f.ValidationErrors
}

func (f *form) GetError(field string) string {
	if f.ValidationErrors == nil {
		return ""
	}

	return f.ValidationErrors.get(field)
}

/*func (f *form) ValidateNonEmpty(field string) {
	if strings.TrimSpace(field) == "" {
		f.ValidationErrors.Add(field, "This field cannot be blank")
	}
}*/

func (f *form) ValidateField(field string, message string, rule ...validation.Rule) *form {
	// If f.Error has an unexpected error, just return form and do nothing.
	// If en unexpected error is present, calling code will ignore validation
	// errors anyways
	if f.unexpectedErr != nil {
		return f
	}

	err := validation.Value(f.Values.Get(field), rule...)
	if err != nil {
		if _, ok := err.(validation.Error); !ok {
			// save unexpected error if empty, do nothing otherwise
			if f.unexpectedErr == nil {
				f.unexpectedErr = err
			}
		} else {
			if f.ValidationErrors == nil {
				f.ValidationErrors = &validationError{}
			}

			var msg string
			if message == "" {
				msg = message
			} else {
				msg = err.Error()
			}
			f.ValidationErrors.Add(field, msg)
		}
	}

	return f
}
