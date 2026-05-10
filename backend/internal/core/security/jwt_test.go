package security

import (
	"context"
	"errors"
	"testing"
)

func TestGenerateAndParseJWT(t *testing.T) {
	token, err := GenerateJWT("secret", "user-id", "admin")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	claims, err := ParseJWT("secret", token)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if claims.UserID != "user-id" {
		t.Fatalf("expected user-id, got %s", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Fatalf("expected admin, got %s", claims.Role)
	}
}

func TestParseJWTInvalid(t *testing.T) {
	token, err := GenerateJWT("secret", "user-id", "admin")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	_, err = ParseJWT("other-secret", token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestClaimsContext(t *testing.T) {
	ctx := WithClaims(context.Background(), Claims{
		UserID: "user-id",
		Role:   RoleUser,
	})

	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		t.Fatalf("expected claims in context")
	}
	if claims.UserID != "user-id" || claims.Role != RoleUser {
		t.Fatalf("unexpected claims: %#v", claims)
	}

	_, ok = ClaimsFromContext(context.Background())
	if ok {
		t.Fatalf("expected no claims in empty context")
	}
}
