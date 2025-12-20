package authService

import "testing"

func TestValidateEmail_Empty(t *testing.T) {
	if validateEmail("") {
		t.Fatalf("expected false for empty email")
	}
}
