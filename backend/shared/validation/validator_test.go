package validation

import (
	"strings"
	"testing"
)

func TestValidator_Required(t *testing.T) {
	v := New()
	v.Required("name", "")
	if err := v.Validate(); err == nil {
		t.Error("empty required field should fail")
	}

	v2 := New()
	v2.Required("name", "John")
	if err := v2.Validate(); err != nil {
		t.Errorf("non-empty required field should pass: %v", err)
	}
}

func TestValidator_Email(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"user.name+tag@example.co.kr", true},
		{"invalid", false},
		{"@example.com", false},
		{"user@", false},
		{"", true}, // empty is valid (use Required for non-empty check)
	}

	for _, tt := range tests {
		v := New()
		v.Email("email", tt.email)
		err := v.Validate()
		if tt.valid && err != nil {
			t.Errorf("email %q should be valid: %v", tt.email, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("email %q should be invalid", tt.email)
		}
	}
}

func TestValidator_UUID(t *testing.T) {
	v := New()
	v.UUID("id", "550e8400-e29b-41d4-a716-446655440000")
	if err := v.Validate(); err != nil {
		t.Errorf("valid UUID should pass: %v", err)
	}

	v2 := New()
	v2.UUID("id", "not-a-uuid")
	if err := v2.Validate(); err == nil {
		t.Error("invalid UUID should fail")
	}
}

func TestValidator_MinMaxLength(t *testing.T) {
	v := New()
	v.MinLength("password", "ab", 8)
	if err := v.Validate(); err == nil {
		t.Error("short string should fail MinLength")
	}

	v2 := New()
	v2.MaxLength("name", "a very long name that exceeds the limit", 10)
	if err := v2.Validate(); err == nil {
		t.Error("long string should fail MaxLength")
	}
}

func TestValidator_Range(t *testing.T) {
	v := New()
	v.Range("score", 150.0, 0.0, 100.0)
	if err := v.Validate(); err == nil {
		t.Error("out of range value should fail")
	}
}

func TestValidator_OneOf(t *testing.T) {
	v := New()
	v.OneOf("role", "superuser", []string{"admin", "user", "guest"})
	if err := v.Validate(); err == nil {
		t.Error("value not in allowed set should fail")
	}
}

func TestValidator_Chaining(t *testing.T) {
	v := New()
	v.Required("email", "").
		Email("email", "invalid").
		MinLength("password", "ab", 8)

	errs := v.Errors()
	if len(errs) != 3 {
		t.Errorf("expected 3 errors, got %d", len(errs))
	}
}

func TestSanitizeString(t *testing.T) {
	input := "  <script>alert('xss')</script>  "
	result := SanitizeString(input)
	if strings.Contains(result, "<script>") {
		t.Error("HTML should be escaped")
	}
	if strings.HasPrefix(result, " ") || strings.HasSuffix(result, " ") {
		t.Error("whitespace should be trimmed")
	}
}
