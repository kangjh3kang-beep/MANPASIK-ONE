package validation

import (
	"html"
	"strings"
)

// SanitizeString removes potentially dangerous characters from user input
func SanitizeString(s string) string {
	// Trim whitespace
	s = strings.TrimSpace(s)
	// Escape HTML entities to prevent XSS
	s = html.EscapeString(s)
	return s
}

// SanitizeMultiline removes dangerous chars but preserves newlines
func SanitizeMultiline(s string) string {
	s = strings.TrimSpace(s)
	s = html.EscapeString(s)
	return s
}
