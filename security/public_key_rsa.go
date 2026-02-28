package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GetPublicKeyFromFile(location string) (publicKey *rsa.PublicKey, err error) {
	// read the PEM file
	byteData, errRead := os.ReadFile(location)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read rsa public key file Error: %w", errRead)
	}

	return GetPublicKey(byteData)
}

func GetPublicKey[T []byte | string](pemData T) (publicKey *rsa.PublicKey, err error) {
	// convert to []byte for decoding
	var data []byte
	switch v := any(pemData).(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil, fmt.Errorf("unsupported type %T", pemData)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("pem data cannot be empty")
	}

	// decode the PEM block
	blockData, _ := pem.Decode(data)
	if blockData == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// parse the RSA public key
	publicKeyData, errParse := x509.ParsePKIXPublicKey(blockData.Bytes)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse rsa public key Error: %w", errParse)
	}

	// assert the type to *rsa.PublicKey
	publicKey, ok := publicKeyData.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type RSA")
	}
	return publicKey, nil
}
