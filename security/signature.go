package security

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func GenerateAsymmetricSignature(privateKey *rsa.PrivateKey, signToString []byte, hashType crypto.Hash) (string, error) {
	hashed, errHash := Digest(hashType, signToString)
	if errHash != nil {
		return "", fmt.Errorf("failed to generate hash Error: %w", errHash)
	}

	signature, errSign := rsa.SignPKCS1v15(rand.Reader, privateKey, hashType, hashed)
	if errSign != nil {
		return "", fmt.Errorf("failed to generate signature Error: %w", errSign)
	}

	return HmacShaEncode(signature), nil
}

func GenerateSymmetricSignature(secretKey string, signToString []byte, hashType crypto.Hash) string {

	hasher := hmac.New(hashType.New, []byte(secretKey))
	hasher.Write(signToString)
	signature := hasher.Sum(nil)

	return HmacShaEncode(signature)

}
