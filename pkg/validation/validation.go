package validation

import (
	"reflect"
	"time"
)

// Rule is the interface all validation rule must implement
type Rule interface {
	Validate(v interface{}) error
}

func Value(v interface{}, rules ...Rule) error {
	var isRequired bool
	for _, r := range rules {
		if reflect.ValueOf(r).Pointer() == reflect.ValueOf(Required).Pointer() {
			isRequired = true
			break
		}
	}

	// if the required rule is absent and field being validated is empty
	// skip validation
	if !isRequired && isEmpty(v) {
		return nil
	}

	for _, r := range rules {
		err := r.Validate(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func isEmpty(field interface{}) bool {
	switch f := field.(type) {
	case string:
		if f == "" {
			return true
		}
	case int, int32, int64:
		if f == 0 {
			return true
		}
	case time.Time:
		if f.IsZero() {
			return true
		}
	case *time.Time:
		if f == nil {
			return true
		}
	}

	return false
}

// The RuleFunc type is an adapter to allow the use of ordinary functions as Rules.
// If f is a function with the appropriate signature, RuleFunc(f) is a Rule that calls f.
// validation.Rule
type RuleFunc func(interface{}) error

func (f RuleFunc) Validate(v interface{}) error {
	return f(v)
}
