package fluent_validator

import (
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	isNonEmpty := func(s string) ValidatorFunc {
		return func() ValidationResult {
			if s == "" {
				return ValidationResult{IsValid: false, Message: []string{"must not be empty"}}
			}
			return ValidationResult{IsValid: true}
		}
	}
	hasMinLen := func(s string, n int) ValidatorFunc {
		return func() ValidationResult {
			if len(s) < n {
				return ValidationResult{IsValid: false, Message: []string{"too short"}}
			}
			return ValidationResult{IsValid: true}
		}
	}

	tests := []struct {
		name        string
		build       func() *FluentValidator
		wantValid   bool
		wantMessage []string
	}{
		{
			name:        "empty chain is valid",
			build:       func() *FluentValidator { return New() },
			wantValid:   true,
			wantMessage: []string{},
		},
		{
			name: "AND: both pass",
			build: func() *FluentValidator {
				return New().
					And(isNonEmpty("abc")).
					And(hasMinLen("abc", 3))
			},
			wantValid:   true,
			wantMessage: []string{},
		},
		{
			name: "AND: first fails",
			build: func() *FluentValidator {
				return New().
					And(isNonEmpty("")).
					And(hasMinLen("abc", 3))
			},
			wantValid:   false,
			wantMessage: []string{"must not be empty"},
		},
		{
			name: "OR: one passes",
			build: func() *FluentValidator {
				return New().
					Or(isNonEmpty("")).
					Or(isNonEmpty("x"))
			},
			wantValid: true,
		},
		{
			name: "OR: both fail",
			build: func() *FluentValidator {
				return New().
					Or(isNonEmpty("")).
					Or(hasMinLen("a", 2))
			},
			wantValid:   false,
			wantMessage: []string{"must not be empty", "too short"},
		},
		{
			name: "AND: multiple failures only collect first (short-circuit)",
			build: func() *FluentValidator {
				failMsg := func(msg string) ValidatorFunc {
					return func() ValidationResult { return Fail(msg) }
				}
				return New().
					And(failMsg("e1")).
					And(failMsg("e2")).
					And(failMsg("e3"))
			},
			wantValid:   false,
			wantMessage: []string{"e1"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := tc.build().Validate()
			if res.IsValid != tc.wantValid {
				t.Fatalf("%s: expected valid=%v, got %v", tc.name, tc.wantValid, res.IsValid)
			}
			if tc.wantMessage != nil && !reflect.DeepEqual(res.Message, tc.wantMessage) {
				t.Fatalf("%s: expected messages=%v, got %v", tc.name, tc.wantMessage, res.Message)
			}
		})
	}
}
