package main

import "testing"

func TestGetEnv(t *testing.T) {
	if v := getEnv("NONEXISTENT_VAR", "default"); v != "default" {
		t.Errorf("expected 'default', got '%s'", v)
	}

	t.Setenv("TEST_VAR", "custom")
	if v := getEnv("TEST_VAR", "default"); v != "custom" {
		t.Errorf("expected 'custom', got '%s'", v)
	}
}
