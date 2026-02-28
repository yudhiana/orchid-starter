package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
)

func generateTestPrivateKeyPEM(t *testing.T) string {
	t.Helper() // Mark this as a test helper

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	pkcs8Key, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Key,
	}

	return string(pem.EncodeToMemory(pemBlock))
}

func TestGetPrivateKey_WithByteSlice(t *testing.T) {
	pemData := generateTestPrivateKeyPEM(t)

	privateKey, err := GetPrivateKey(pemData)

	if err != nil {
		t.Fatalf("GetPrivateKey failed: %v", err)
	}
	if privateKey == nil {
		t.Fatal("expected non-nil private key")
	}
	if privateKey.N == nil {
		t.Fatal("expected valid RSA key with modulus")
	}
}

func TestGetPrivateKey_WithString(t *testing.T) {
	pemData := generateTestPrivateKeyPEM(t)

	privateKey, err := GetPrivateKey(pemData)

	if err != nil {
		t.Fatalf("GetPrivateKey failed: %v", err)
	}
	if privateKey == nil {
		t.Fatal("expected non-nil private key")
	}
	if privateKey.N == nil {
		t.Fatal("expected valid RSA key with modulus")
	}
}

func TestGetPrivateKey_EmptyData(t *testing.T) {
	_, err := GetPrivateKey("")

	if err == nil {
		t.Fatal("expected error for empty data")
	}
}

func TestGetPrivateKey_InvalidPEMFormat(t *testing.T) {
	invalidPEM := "this is not valid PEM data"

	_, err := GetPrivateKey(invalidPEM)

	if err == nil {
		t.Fatal("expected error for invalid PEM format")
	}
}

func TestGetPrivateKey_InvalidKeyData(t *testing.T) {
	invalidKeyData := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: []byte("invalid key data"),
	}

	invalidPEM := string(pem.EncodeToMemory(invalidKeyData))

	_, err := GetPrivateKey(invalidPEM)

	if err == nil {
		t.Fatal("expected error for invalid key data")
	}
}

func TestGetPrivateKeyFromFile_Success(t *testing.T) {
	pemData := generateTestPrivateKeyPEM(t)

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test-private-key-*.pem")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write PEM data to file
	_, err = tmpFile.WriteString(pemData)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test reading from file
	privateKey, err := GetPrivateKeyFromFile(tmpFile.Name())

	if err != nil {
		t.Fatalf("GetPrivateKeyFromFile failed: %v", err)
	}
	if privateKey == nil {
		t.Fatal("expected non-nil private key")
	}
}

func TestGetPrivateKeyFromFile_FileNotFound(t *testing.T) {
	_, err := GetPrivateKeyFromFile("/nonexistent/path/to/key.pem")

	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestGetPrivateKey_CorrectKeyType(t *testing.T) {
	pemData := generateTestPrivateKeyPEM(t)

	privateKey, err := GetPrivateKey(pemData)

	if err != nil {
		t.Fatalf("GetPrivateKey failed: %v", err)
	}

	// Verify it's an RSA key by checking it can be used for signing
	testData := []byte("test")
	_, err = GenerateAsymmetricSignature(privateKey, testData, 5) // crypto.SHA256 = 5

	if err != nil {
		t.Fatalf("private key cannot be used for signing: %v", err)
	}
}
