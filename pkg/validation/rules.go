package validation

import (
	"errors"
	"fmt"
	validation2 "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"unicode/utf8"
)

// Required is a rule for values that must be non-zero
// Validates only string, int, and float values
var Required = RuleFunc(
	func(v interface{}) error {
		switch resolved := v.(type) {
		case string:
			if resolved == "" {
				return NewError("string is empty")
			}
		case float64, float32, int64, int32, int16, int8, int, uint:
			if resolved == 0 {
				return NewError("float is empty")
			}
		default:
			return errors.New("unexpected type")
		}

		return nil
	},
)

/*type length struct {
	min int
	max int
}
func (l length) Validate(v interface{}) error {
	err := validation2.Validate(v, validation2.Length(l.min, l.max))
	if err != nil {
		if _, ok := err.(validation2.InternalError); !ok {
			return NewError("field is required")
		}
		return err
	}
	return nil
}

// Length returns a validation rule that checks if a value's length is within the specified range.
func Length(min int, max int) Rule {
	return length{min:min, max:max}
}*/

var Length = func(min, max int) Rule {
	return RuleFunc(
		func(v interface{}) error {
			str, ok := v.(string)
			if !ok {
				return errors.New("value is not a string")
			}
			n := utf8.RuneCountInString(str)
			if n < min {
				return NewError(fmt.Sprintf("value(%d) cannot be less than min(%d)", n, min))
			}
			if n > max {
				return NewError(fmt.Sprintf("value(%d) cannot be greater than max(%d)", n, min))
			}
			return nil
		},
	)
}

var Email = RuleFunc(
	func(v interface{}) error {
		if s, ok := v.(string); !ok && len(s) > 128 {
			return NewError("email too long")
		}

		err := validation2.Validate(v, is.Email)
		if err != nil {
			if _, ok := err.(validation2.InternalError); !ok {
				return NewError("not a valid email")
			}
			return err
		}
		return nil
	},
)

/*// phone
type phone struct {}
func (e phone) Validate(v interface{}) error {
	if s, ok := v.(string); !ok && len(s) > 32 {
		return errors.New("phone too long")
	}

	err := validation2.Validate(v, is.E164)
	if err != nil {
		return errors.New("not a phone number")
	}

	return nil
}

var Phone = phone{}*/

/*// InRange returns a validation rule that checks if an integer is within the specified range.
func InRange(min, max int) Rule {
	return inRange{min:min, max:max}
}

type inRange struct {
	min int
	max int
}

func (r inRange) Validate(v interface{}) error {
	num, ok := v.(int)
	if !ok {
		return fmt.Errorf("value is not between %d and %d inclusive", r.min, r.max)
	}
	if num < r.min || num > r.max {
		return fmt.Errorf("value is not between %d and %d inclusive", r.min, r.max)
	}
	return nil
}*/

// InRange is a function that returns a validation rule that checks if an integer
// is within the specified range.
var InRange = func(min, max int) Rule { // todo:: max = 0 should allow value of any size
	return RuleFunc(
		func(v interface{}) error {
			num, ok := v.(int)
			if !ok {
				return errors.New("could not convert value to int")
			}
			if num < min || num > max {
				return NewError(fmt.Sprintf("value is not between %d and %d inclusively", min, max))
			}
			return nil
		},
	)
}

/*// single line
type singleLineValidator struct{}

func (*singleLineValidator) Validate(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("invalid string")
	}

	if strings.Contains(str, "\n") {
		return errors.New("field cannot contain new line character")
	}

	return nil
}

var SingleLineRule = &singleLineValidator{}*/
