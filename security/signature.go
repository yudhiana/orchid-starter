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

	return EncodeBase64Key(signature), nil
}

func GenerateSymmetricSignature(secretKey string, signToString []byte, hashType crypto.Hash) string {

	hasher := hmac.New(hashType.New, []byte(secretKey))
	hasher.Write(signToString)
	signature := hasher.Sum(nil)

	return EncodeBase64Key(signature)
}

func VerifySymmetricSignature(secretKey string, signToString []byte, incomingSignature string, hashType crypto.Hash) (bool, error) {
	// Regenerate the expected signature
	expectedSign := GenerateSymmetricSignature(secretKey, signToString, hashType)

	// Constant time comparison prevents timing attacks
	isValid := hmac.Equal([]byte(expectedSign), []byte(incomingSignature))

	return isValid, nil
}

func VerifyAsymmetricSignature(publicKey *rsa.PublicKey, signToString []byte, incomingSignature string, hashType crypto.Hash) (bool, error) {
	// Decode the base64-encoded incoming signature
	signatureBytes, err := DecodeBase64Key(incomingSignature)
	if err != nil {
		return false, fmt.Errorf("failed to decode incoming signature: %w", err)
	}

	// Hash the original data using the same algorithm
	hashed, err := Digest(hashType, signToString)
	if err != nil {
		return false, fmt.Errorf("failed to generate hash: %w", err)
	}

	// Verify the signature using the public key
	err = rsa.VerifyPKCS1v15(publicKey, hashType, hashed, signatureBytes)
	if err != nil {
		// Signature is invalid but not an error condition
		return false, nil
	}

	return true, nil
}
