package security

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func makeBase64Key(t *testing.T) string {
	t.Helper()
	raw := []byte("test-secret-key")
	return base64.StdEncoding.EncodeToString(raw)
}

func TestGenerateJwtToken(t *testing.T) {
	secret := makeBase64Key(t)

	claim := &AppClaim{
		Username: "alice",
		Email:    "alice@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "unit-test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}

	tok, err := claim.GenerateJwtToken(jwt.SigningMethodHS256, secret, WithHeader("kid", "1234"))
	if err != nil {
		t.Fatalf("GenerateJwtToken error: %v", err)
	}
	if tok == "" {
		t.Fatal("expected non-empty token string")
	}
}
