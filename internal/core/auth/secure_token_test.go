package core_auth_test

import (
	"testing"

	core_auth "github.com/daf32/golang-todoapp/internal/core/auth"
)

func TestGenerateSecuteToken_UniqueAndNonEmpty(t *testing.T) {
	a, err := core_auth.GenerateSecureToken(32)
	if err != nil {
		t.Fatal(err)
	}

	b, err := core_auth.GenerateSecureToken(32)
	if err != nil {
		t.Fatal(err)
	}

	if a == "" || a == b {
		t.Fatalf("expected unique non-empty tokens, got %q and %q", a, b)
	}
}
