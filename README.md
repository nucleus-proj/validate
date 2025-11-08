# Fluent Validator

A tiny fluent validation helper for composing validation rules with AND/OR semantics. It evaluates validators left-to-right, short-circuits logically where possible, and returns a `ValidationResult` including optional failure messages.

## Features

- Zero config; tiny API surface
- AND/OR composition with left-to-right evaluation and short-circuiting
- Aggregated failure messages (AND collects failing step messages; OR clears when any passes)
- Simple function adapter via `ValidatorFunc`
- Bring-your-own rules: implement any `Validator`; no built-in rule set required
- Built-in common rules (string/number/time/collection)

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"core/pkg/fluent_validator"
)

func main() {
	// Validate an email address
	email := "user@example.com"
	
	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(email)).
		And(fluent_validator.EmailValid(email)).
		Validate()

	if !result.IsValid {
		log.Fatalf("Validation failed: %v", result.Message)
	}

	fmt.Println("Email is valid!")

	// Validate a password with multiple rules
	password := "SecurePass123"
	
	result = fluent_validator.New().
		And(fluent_validator.NonEmpty(password)).
		And(fluent_validator.MinLen(password, 8)).
		And(fluent_validator.MaxLen(password, 50)).
		Validate()

	if !result.IsValid {
		log.Fatalf("Validation failed: %v", result.Message)
	}

	fmt.Println("Password is valid!")

	// Validate a number range
	age := 25
	
	result = fluent_validator.New().
		And(fluent_validator.IntMin(age, 18)).
		And(fluent_validator.IntMax(age, 120)).
		Validate()

	if !result.IsValid {
		log.Fatalf("Validation failed: %v", result.Message)
	}

	fmt.Println("Age is valid!")
}
```

## Versioning

This package follows Semantic Versioning (SemVer):

- MAJOR (X.0.0): breaking API changes
- MINOR (0.X.0): backward-compatible features and improvements
- PATCH (0.0.X): backward-compatible bug fixes

## Upcoming

1.x (minor, non-breaking): continue adding rules, helpers, and docs
  - JSON-driven test generation helper (`gen-fv-tests`)
  - Field-level helpers for struct validation (tags or small DSL)
  - Result metadata (code/field) for UI/form integration
  - Configurable evaluation mode (short-circuit vs exhaustive; defaults unchanged)
  - Message templating/localization (placeholders, i18n)
- 2.0 (major, breaking):
    - Grouping/precedence semantics if they alter evaluation outcomes
    - Context-aware validators if `Validate()` signature changes


## Usage

```go
import (
	"core/pkg/fluent_validator"
)

// Define some rule functions
isNonEmpty := func(s string) fluent_validator.ValidatorFunc {
	return func() fluent_validator.ValidationResult {
		if s == "" {
            return fluent_validator.ValidationResult{IsValid: false, Message: []string{"must not be empty"}}
		}
		return fluent_validator.ValidationResult{IsValid: true}
	}
}

hasMinLen := func(s string, n int) fluent_validator.ValidatorFunc {
	return func() fluent_validator.ValidationResult {
		if len(s) < n {
            return fluent_validator.ValidationResult{IsValid: false, Message: []string{"too short"}}
		}
		return fluent_validator.ValidationResult{IsValid: true}
	}
}

// Compose with AND/OR
v := fluent_validator.New().
	And(isNonEmpty("abc")).
	And(hasMinLen("abc", 3)).
	Or(fluent_validator.ValidatorFunc(func() fluent_validator.ValidationResult {
		return fluent_validator.ValidationResult{IsValid: true}
	}))

res := v.Validate()
if !res.IsValid {
    // handle res.Message (slice of strings)
}
```

## API

- `type ValidationResult struct { IsValid bool; Message []string }`
- `type Validator interface { Validate() ValidationResult }`
- `type ValidatorFunc func() ValidationResult`
- `func Success() ValidationResult`
- `func Fail(msg ...string) ValidationResult`
- `func New() *FluentValidator`
- `func (*FluentValidator) And(v Validator) *FluentValidator`
- `func (*FluentValidator) Or(v Validator) *FluentValidator`
- `func (*FluentValidator) Validate() ValidationResult`

Built-in rules:
- String: `NonEmpty`, `MinLen`, `MaxLen`, `LenBetween`, `Matches`, `OneOf`
- Number: `IntMin`, `IntMax`, `IntBetween`, `IntNonZero`, `FloatMin`, `FloatMax`, `FloatBetween`, `FloatNonZero`
- Time: `TimeNotZero`, `TimeBefore`, `TimeAfter`, `TimeBetween`
- Collection: `NotEmptyLen`, `LenMin`, `LenMax`, `LenBetweenSize`, `ContainsString`, `UniqueStrings`
- Contact: `EmailValid`, `PhoneE164`
### Notes

// Evaluation is left-to-right and short-circuits within contiguous AND/OR segments:
- AND: requires all validators to pass; collects failures up to and including the first failure
- OR: passes if any validator passes; collects all failures only if all fail, and clears messages when any passes
