package validation

import (
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	tests := []struct {
		name string
		value interface{}
		rule Rule
		expectedErrType ErrType
	}{
		{
			name: "valid string | required",
			value: "hello",
			rule: Required,
			expectedErrType: errTypeNil,
		},
		{
			name: "empty string | required",
			value: "",
			rule: Required,
			expectedErrType: errTypeValidation,
		},
		{
			name: "unexpected value | required",
			value: struct{}{},
			rule: Required,
			expectedErrType: errTypeUnexpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Value(tt.value, tt.rule)
			switch tt.expectedErrType {
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
