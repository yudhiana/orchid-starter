package security

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

type TokenOption func(*jwt.Token)

func WithHeader(key string, value any) TokenOption {
	return func(t *jwt.Token) {
		t.Header[key] = value
	}
}

type AppClaim struct {
	jwt.RegisteredClaims
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (c *AppClaim) GenerateJwtToken(method jwt.SigningMethod, secretKey string, options ...TokenOption) (token string, err error) {
	jwtToken := jwt.NewWithClaims(method, c)
	for _, option := range options {
		option(jwtToken)
	}

	secretKeyBytes, errDecode := DecodeBase64Key(secretKey)
	if errDecode != nil {
		err = fmt.Errorf("failed to decode secret key Error: %w", errDecode)
		return
	}

	return jwtToken.SignedString(secretKeyBytes)
}

func ValidateJWTToken(tokenJWT string, secretKey string, method jwt.SigningMethod) (claim *AppClaim, valid bool, err error) {
	// use ParseWithClaims to ensure the custom claim type is populated
	appClaim := new(AppClaim)
	jwtToken, errJwt := jwt.ParseWithClaims(tokenJWT, appClaim, func(token *jwt.Token) (any, error) {
		if signMethod, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		} else if method != nil && signMethod != method {
			return nil, ErrInvalidSigningMethod
		}
		return DecodeBase64Key(secretKey)
	})

	if errJwt != nil {
		err = fmt.Errorf("failed to validate token Error: %w", errJwt)
		return
	}

	if jwtToken != nil && jwtToken.Valid {
		// claims have been written to appClaim by ParseWithClaims
		claim = appClaim
		valid = true
	} else {
		valid = false
	}
	return
}
