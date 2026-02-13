package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ValidationError holds validation error details
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []*ValidationError

func (e ValidationErrors) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "; ")
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Validator provides fluent validation
type Validator struct {
	errors ValidationErrors
}

// New creates a new Validator
func New() *Validator {
	return &Validator{}
}

// Validate returns errors if any
func (v *Validator) Validate() error {
	if v.errors.HasErrors() {
		return v.errors
	}
	return nil
}

// Errors returns the validation errors
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

func (v *Validator) addError(field, message string) *Validator {
	v.errors = append(v.errors, &ValidationError{Field: field, Message: message})
	return v
}

// Required checks that a string field is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.addError(field, "필수 입력 항목입니다")
	}
	return v
}

// MinLength checks minimum string length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if utf8.RuneCountInString(value) < min {
		v.addError(field, fmt.Sprintf("최소 %d자 이상이어야 합니다", min))
	}
	return v
}

// MaxLength checks maximum string length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if utf8.RuneCountInString(value) > max {
		v.addError(field, fmt.Sprintf("최대 %d자까지 허용됩니다", max))
	}
	return v
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Email validates email format
func (v *Validator) Email(field, value string) *Validator {
	if value != "" && !emailRegex.MatchString(value) {
		v.addError(field, "유효한 이메일 형식이 아닙니다")
	}
	return v
}

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// UUID validates UUID format
func (v *Validator) UUID(field, value string) *Validator {
	if value != "" && !uuidRegex.MatchString(value) {
		v.addError(field, "유효한 UUID 형식이 아닙니다")
	}
	return v
}

var phoneRegex = regexp.MustCompile(`^(\+?\d{1,4}[-.\s]?)?\(?\d{1,4}\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}$`)

// Phone validates phone number format
func (v *Validator) Phone(field, value string) *Validator {
	if value != "" && !phoneRegex.MatchString(value) {
		v.addError(field, "유효한 전화번호 형식이 아닙니다")
	}
	return v
}

// Range checks numeric range (float64)
func (v *Validator) Range(field string, value, min, max float64) *Validator {
	if value < min || value > max {
		v.addError(field, fmt.Sprintf("%.2f ~ %.2f 범위여야 합니다", min, max))
	}
	return v
}

// PositiveInt checks that an integer is positive
func (v *Validator) PositiveInt(field string, value int64) *Validator {
	if value <= 0 {
		v.addError(field, "양수여야 합니다")
	}
	return v
}

// OneOf checks that a value is in a set of allowed values
func (v *Validator) OneOf(field, value string, allowed []string) *Validator {
	if value == "" {
		return v
	}
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.addError(field, fmt.Sprintf("허용된 값: %s", strings.Join(allowed, ", ")))
	return v
}
