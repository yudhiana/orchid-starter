package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GetRSAFromFile(location string) (privateKey *rsa.PrivateKey, err error) {
	// read the PEM file
	byteData, errRead := os.ReadFile(location)
	if errRead != nil {
		err = fmt.Errorf("failed to read rsa file Error: %w", errRead)
	}

	// decode the PEM block
	blockData, _ := pem.Decode(byteData)

	// parse the RSA private key
	privateKeyData, errParse := x509.ParsePKCS8PrivateKey(blockData.Bytes)
	if errParse != nil {
		err = fmt.Errorf("failed to parse rsa private key Error: %w", errParse)
		return
	}

	// assert the type to *rsa.PrivateKey
	privateKey, ok := privateKeyData.(*rsa.PrivateKey)
	if !ok {
		err = fmt.Errorf("private key is not of type RSA")
		return
	}
	return
}
