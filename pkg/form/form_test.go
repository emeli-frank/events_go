package form

import (
	"events/pkg/validation"
	"net/url"
	"testing"
)

func TestForm_ValidateField(t *testing.T) {
	tests := []struct{
		name                   string
		value                  url.Values
		expectedNumberOfErrors int
	} {
		{
			name: "field with one valid string",
			value: map[string][]string{
				"name": {"Samwise"},
			},
			expectedNumberOfErrors: 0,
		},
		{
			name: "field with one invalid string",
			value: map[string][]string{
				"name": {""},
			},
			expectedNumberOfErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.value).ValidateField("name", validation.Required).Error()

			if err == nil {
				if tt.expectedNumberOfErrors != 0 {
					t.Fatalf("expected %d error(s), got nil", tt.expectedNumberOfErrors)
				}
			} else {
				e, ok := err.(Error)
				if !ok {
					t.Fatalf("a validation error was not received, error received is: %v", err)
				}
				m := e.ErrorMessages()
				if numberOfErrors(m) != tt.expectedNumberOfErrors {
					t.Fatalf("expected %d error(s), got %d", tt.expectedNumberOfErrors, numberOfErrors(m))
				}
			}


			/*if err == nil && tt.expectedNumberOfErrors != 0 {
				t.Fatalf("expected %d error(s), got nil", tt.expectedNumberOfErrors)
			}
			else {
				e, ok := err.(Error)
				if !ok {
					t.Fatalf("a validation error was not received, error received is: %v", err)
				}
				m := e.ErrorMessages()
				if numberOfErrors(m) != tt.expectedNumberOfErrors {
					t.Fatalf("expected %d error(s), got %d", tt.expectedNumberOfErrors, numberOfErrors(m))
				}
			}*/

			/*e, ok := err.(Error)
			if !ok {
				t.Fatal("a validation error was not received")
			}
			m := e.ErrorMessages()
			if numberOfErrors(m) != tt.expectedNumberOfErrors {
				t.Fatalf("expected %d error(s), got %d", tt.expectedNumberOfErrors, numberOfErrors(m))
			}*/
		})
	}
}

func numberOfErrors(m map[string][]string) int {
	count := 0
	for _, v := range m {
		count += len(v)
	}
	return count
}
