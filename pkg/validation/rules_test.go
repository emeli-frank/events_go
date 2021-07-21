package validation

import (
	"reflect"
	"testing"
)

var (
	//ErrTypeValidation = reflect.TypeOf(NewError(""))
)

type ErrType int
const (
	errTypeValidation  ErrType= 1
	errTypeUnexpected ErrType = 2
	errTypeNil ErrType = 3
)

func TestRules(t *testing.T) {
	tests := []struct{
		name string
		value interface{}
		rule Rule
		wantedErrType ErrType
	} {
		{
			name: "Non empty string",
			value: "Hello",
			rule: Required,
			wantedErrType: errTypeNil,
		},
		{
			name: "Empty string",
			value: "",
			rule: Required,
			wantedErrType: errTypeValidation,
		},
		{
			name: "Non-zero int",
			value: 54,
			rule: Required,
			wantedErrType: errTypeNil,
		},
		{
			name: "Zero int",
			value: 0,
			rule: Required,
			wantedErrType: errTypeValidation,
		},

		{
			name: "Unexpected type | required",
			value: struct {}{},
			rule: Required,
			wantedErrType: errTypeUnexpected,
		},
		{
			name: "Valid email",
			value: "email@example.com",
			rule: Email,
			wantedErrType: errTypeNil,
		},
		{
			name: "Short valid email",
			value: "e@e.c",
			rule: Email,
			wantedErrType: errTypeNil,
		},
		{
			name: "Invalid email: missing @",
			value: "emailexample.com",
			rule: Email,
			wantedErrType: errTypeValidation,
		},
		{
			name: "Invalid email: missing tld",
			value: "email@examplecom",
			rule: Email,
			wantedErrType: errTypeValidation,
		},
		{
			name: "String length withing range",
			value: "cat",
			rule: Length(2, 5),
			wantedErrType: errTypeNil,
		},
		{
			name: "Value equal to min length value",
			value: "do",
			rule: Length(2, 5),
			wantedErrType: errTypeNil,
		},
		{
			name: "Value equal to max length value",
			value: "hello",
			rule: Length(2, 5),
			wantedErrType: errTypeNil,
		},
		{
			name: "Value shorter than min",
			value: "i",
			rule: Length(2, 5),
			wantedErrType: errTypeValidation,
		},
		{
			name: "Value longer than max",
			value: "golang",
			rule: Length(2, 5),
			wantedErrType: errTypeValidation,
		},

		{
			name: "Unexpected type | Length",
			value: struct {}{},
			rule: Length(2, 5),
			wantedErrType: errTypeUnexpected,
		},
		{
			name: "Int withing range",
			value: 3,
			rule: InRange(2, 5),
			wantedErrType: errTypeNil,
		},
		{
			name: "Value equal to min range value",
			value: 2,
			rule: InRange(2, 5),
			wantedErrType: errTypeNil,
		},
		{
			name: "Value equal to max range value",
			value: 5,
			rule: InRange(2, 5),
			wantedErrType: errTypeNil,
		},
		{
			name: "Value lower than min",
			value: 1,
			rule: InRange(2, 5),
			wantedErrType: errTypeValidation,
		},
		{
			name: "Value higher than max",
			value: 6,
			rule: InRange(2, 5),
			wantedErrType: errTypeValidation,
		},
		{
			name: "Unexpected type | InRange",
			value: struct {}{},
			rule: InRange(2, 5),
			wantedErrType: errTypeUnexpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate(tt.value)

			switch tt.wantedErrType {
			case errTypeNil:
				if err != nil {
					t.Fatalf("wanted nil, got error of type: %T and value: %v", err, err)
				}
			case errTypeValidation:
				if reflect.TypeOf(err).String() != reflect.TypeOf(NewError("")).String() {
					t.Fatalf("wanted validation error of type %v, got type: %T and value: %v",
						reflect.TypeOf(NewError("")).String(), err, err)
				}
			case errTypeUnexpected:
				if reflect.TypeOf(err).String() == reflect.TypeOf(NewError("")).String() {
					t.Fatalf("wanted unexpected error, got type: %T and value: %v", err, err)
				}
			}
		})
	}
}
