package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func GenerateSignature(priv *rsa.PrivateKey, signString string) (string, error) {
	hashString, err := GenerateHash(signString, crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash Error: %w", err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashString)
	if err != nil {
		return "", fmt.Errorf("failed to generate signature Error: %w", err)
	}

	return HmacShaEncode(string(signature)), nil
}
