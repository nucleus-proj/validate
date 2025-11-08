// Package fluent_validator provides a small fluent validation helper for
// composing multiple validation steps with AND/OR semantics. It evaluates
// validators left-to-right and short-circuits logically where possible
// while preserving an explanatory failure message when invalid.
package fluent_validator

// ValidationResult represents the outcome of a validation step.
type ValidationResult struct {
	IsValid bool
	Message []string
}

// Validator is the contract for any validation step.
// Implementations should return a ValidationResult indicating success
// and, if invalid, an explanatory message.
type Validator interface {
	Validate() ValidationResult
}

// ValidatorFunc is an adapter to allow the use of ordinary functions as validators.
type ValidatorFunc func() ValidationResult

// Validate calls the underlying function.
func (f ValidatorFunc) Validate() ValidationResult { return f() }

// Success returns a successful ValidationResult with an empty message slice.
func Success() ValidationResult { return ValidationResult{IsValid: true, Message: []string{}} }

// Fail returns a failed ValidationResult with the provided messages.
func Fail(msg ...string) ValidationResult { return ValidationResult{IsValid: false, Message: msg} }

// FluentValidator composes multiple validators using AND/OR operators.
// By default, evaluation is left-to-right; AND requires all to pass,
// OR requires at least one to pass. Evaluation short-circuits within
// contiguous AND/OR segments where possible to avoid wasted work:
//   - AND: on first failure, later AND steps are skipped until an OR
//   - OR: on first success, later OR steps are skipped until an AND
//
// Message policy:
//   - AND: collects failures up to and including the first failure
//   - OR: collects all failures if all fail; clears when any passes
type FluentValidator struct {
	steps []chainedStep
}

// New creates a new FluentValidator instance.
func New() *FluentValidator {
	return &FluentValidator{steps: make([]chainedStep, 0, 4)}
}

type logicalOp uint8

const (
	opAnd logicalOp = iota
	opOr
)

type chainedStep struct {
	validator Validator
	op        logicalOp
}

// And adds a validator combined with AND semantics to the chain and
// returns the same builder for fluent chaining.
func (f *FluentValidator) And(v Validator) *FluentValidator {
	f.steps = append(f.steps, chainedStep{validator: v, op: opAnd})
	return f
}

// Or adds a validator combined with OR semantics to the chain and
// returns the same builder for fluent chaining.
func (f *FluentValidator) Or(v Validator) *FluentValidator {
	f.steps = append(f.steps, chainedStep{validator: v, op: opOr})
	return f
}

// Validate evaluates the chain left-to-right, applying AND/OR semantics.
// It short-circuits where possible and returns a ValidationResult
// indicating overall validity. When invalid, Message aggregates failure
// messages encountered according to the logical operators.
func (f *FluentValidator) Validate() ValidationResult {
	if len(f.steps) == 0 {
		return Success()
	}

	accValid := false
	messages := make([]string, 0, len(f.steps))

	for i, step := range f.steps {
		// Always evaluate the first step to seed accumulator
		if i == 0 {
			res := step.validator.Validate()
			accValid = res.IsValid
			if !res.IsValid && len(res.Message) > 0 {
				messages = append(messages, res.Message...)
			}
			continue
		}

		switch step.op {
		case opAnd:
			// Short-circuit: if already false, AND cannot change the outcome
			if !accValid {
				// Skip evaluation to avoid wasted work and extra messages
				continue
			}
			res := step.validator.Validate()
			if !res.IsValid && len(res.Message) > 0 {
				// AND policy: collect up to and including first failure
				messages = append(messages, res.Message...)
			}
			accValid = accValid && res.IsValid
		case opOr:
			// Short-circuit: if already true, OR cannot change the outcome
			if accValid {
				// Skip evaluation to avoid wasted work
				continue
			}
			res := step.validator.Validate()
			if res.IsValid {
				// OR policy: clear failures when chain becomes valid
				messages = []string{}
			} else if len(res.Message) > 0 {
				// Only collected if still failing overall
				messages = append(messages, res.Message...)
			}
			accValid = accValid || res.IsValid
		}
	}

	if accValid {
		return Success()
	}
	return ValidationResult{IsValid: false, Message: messages}
}
