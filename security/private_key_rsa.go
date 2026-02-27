package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GetPrivateKeyFromFile(location string) (privateKey *rsa.PrivateKey, err error) {
	// read the PEM file
	byteData, errRead := os.ReadFile(location)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read rsa file Error: %w", errRead)
	}

	return GetPrivateKey(byteData)
}

func GetPrivateKey[T []byte | string](pemData T) (privateKey *rsa.PrivateKey, err error) {
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

	// parse the RSA private key
	privateKeyData, errParse := x509.ParsePKCS8PrivateKey(blockData.Bytes)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse rsa private key Error: %w", errParse)
	}

	// assert the type to *rsa.PrivateKey
	privateKey, ok := privateKeyData.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type RSA")
	}
	return privateKey, nil
}
