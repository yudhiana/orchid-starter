package security

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// reuse makeBase64Key from jwt_test.go

func TestValidateJWTToken_Success(t *testing.T) {
	secret := makeBase64Key(t)

	claim := &AppClaim{
		Username: "bob",
		Email:    "bob@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "unit-test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}

	tok, err := claim.GenerateJwtToken(jwt.SigningMethodHS256, secret)
	if err != nil {
		t.Fatalf("setup: unable to generate token: %v", err)
	}

	got, valid, err := ValidateJWTToken(tok, secret, jwt.SigningMethodHS256)
	if err != nil {
		t.Fatalf("ValidateJWTToken error: %v", err)
	}
	if !valid {
		t.Fatal("expected token to be valid")
	}
	if got.Username != claim.Username || got.Email != claim.Email {
		t.Errorf("claim mismatch: got %+v, want %+v", got, claim)
	}
}

func TestValidateJWTToken_InvalidSecret(t *testing.T) {
	secret := makeBase64Key(t)
	claim := new(AppClaim)
	tok, err := claim.GenerateJwtToken(jwt.SigningMethodHS256, secret)
	if err != nil {
		t.Fatalf("setup generate token: %v", err)
	}

	wrong := base64.StdEncoding.EncodeToString([]byte("other-key"))
	_, valid, _ := ValidateJWTToken(tok, wrong, jwt.SigningMethodHS256)
	if valid {
		t.Error("expected invalid with wrong secret")
	}
}

func TestValidateJWTToken_Malformed(t *testing.T) {
	secret := makeBase64Key(t)
	_, valid, err := ValidateJWTToken("this.is.not.valid", secret, jwt.SigningMethodHS256)
	if err == nil {
		t.Error("expected parse error for malformed token")
	}
	if valid {
		t.Error("expected valid=false for malformed token")
	}
}
