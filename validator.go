package weeb

import (
	"regexp"
	"strconv"
)

// ValidationFn represents a check the validator can do again some arbitrary input.
// Here fieldName's main use is for nice human readable error messages.
type ValidationFn func(ctx *Context, fieldName string, value string) []string

// Validator represents an instance of a validator and hold the errors that occurred.
type Validator struct {
	ctx    *Context
	errors []string
}

// NewValidator creates a new validator linked to a specific context. The
// context is required for more advanced validators that need to access the
// database or app configuration.
func NewValidator(ctx *Context) *Validator {
	return &Validator{ctx: ctx, errors: []string{}}
}

// Validate validates part of the input using a set of provided validation functions
func (v *Validator) Validate(fieldName string, value string, validations []ValidationFn) {
	for _, validation := range validations {
		errors := validation(v.ctx, fieldName, value)
		v.errors = append(v.errors, errors...)
	}
}

// Valid returns true if the validator encountered no validation errors
func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

// Errors returns the list of validation errors that occurred
func (v *Validator) Errors() []string {
	return v.errors
}

// ValidatePresence validates the input is a non-empty string
func ValidatePresence() ValidationFn {
	return func(ctx *Context, fieldName string, value string) []string {
		if len(value) == 0 {
			return []string{title(fieldName) + " is a required field"}
		}
		return []string{}
	}
}

// ValidateRegexp validates that the input matches a given regular expression
func ValidateRegexp(r *regexp.Regexp) ValidationFn {
	return func(ctx *Context, fieldName string, value string) []string {
		if !r.Match([]byte(value)) {
			return []string{title(fieldName) + " is not in the valid format"}
		}
		return []string{}
	}
}

// ValidateLength validates the lenght of the input is within certain bounds.
// A min of -1 will be ignored. A max of -1 will be ignored. Min and max are
// inclusive so ValidateLength(1,3) passes [1,2,3] but fails [0,4,...]
func ValidateLength(min int, max int) ValidationFn {
	return func(ctx *Context, fieldName string, value string) []string {
		errors := []string{}
		if min != -1 {
			if len(value) < min {
				errors = append(errors, title(fieldName)+" is shorter than "+strconv.Itoa(min))
			}
		}
		if max != -1 {
			if len(value) > max {
				errors = append(errors, title(fieldName)+" is longer than "+strconv.Itoa(min))
			}
		}
		return errors
	}
}
