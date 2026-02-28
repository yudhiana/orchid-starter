package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
)

func generateTestKeyPairPEM() (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Generate private key PEM
	pkcs8Key, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	privatePemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Key,
	}
	privatePEM := string(pem.EncodeToMemory(privatePemBlock))

	// Generate public key PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	publicPemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPEM := string(pem.EncodeToMemory(publicPemBlock))

	return privatePEM, publicPEM
}

func TestGetPublicKey_WithByteSlice(t *testing.T) {
	_, publicPEM := generateTestKeyPairPEM()

	publicKey, err := GetPublicKey(publicPEM)

	if err != nil {
		t.Fatalf("GetPublicKey failed: %v", err)
	}
	if publicKey == nil {
		t.Fatal("expected non-nil public key")
	}
	if publicKey.N == nil {
		t.Fatal("expected valid RSA key with modulus")
	}
}

func TestGetPublicKey_WithString(t *testing.T) {
	_, publicPEM := generateTestKeyPairPEM()

	publicKey, err := GetPublicKey(publicPEM)

	if err != nil {
		t.Fatalf("GetPublicKey failed: %v", err)
	}
	if publicKey == nil {
		t.Fatal("expected non-nil public key")
	}
	if publicKey.N == nil {
		t.Fatal("expected valid RSA key with modulus")
	}
}

func TestGetPublicKey_EmptyData(t *testing.T) {
	_, err := GetPublicKey("")

	if err == nil {
		t.Fatal("expected error for empty data")
	}
}

func TestGetPublicKey_InvalidPEMFormat(t *testing.T) {
	invalidPEM := "this is not valid PEM data"

	_, err := GetPublicKey(invalidPEM)

	if err == nil {
		t.Fatal("expected error for invalid PEM format")
	}
}

func TestGetPublicKey_InvalidKeyData(t *testing.T) {
	invalidKeyData := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: []byte("invalid key data"),
	}

	invalidPEM := string(pem.EncodeToMemory(invalidKeyData))

	_, err := GetPublicKey(invalidPEM)

	if err == nil {
		t.Fatal("expected error for invalid key data")
	}
}

func TestGetPublicKeyFromFile_Success(t *testing.T) {
	_, publicPEM := generateTestKeyPairPEM()

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test-public-key-*.pem")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write PEM data to file
	_, err = tmpFile.WriteString(publicPEM)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test reading from file
	publicKey, err := GetPublicKeyFromFile(tmpFile.Name())

	if err != nil {
		t.Fatalf("GetPublicKeyFromFile failed: %v", err)
	}
	if publicKey == nil {
		t.Fatal("expected non-nil public key")
	}
}

func TestGetPublicKeyFromFile_FileNotFound(t *testing.T) {
	_, err := GetPublicKeyFromFile("/nonexistent/path/to/key.pem")

	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestGetPublicKey_CorrectKeyType(t *testing.T) {
	privatePEM, publicPEM := generateTestKeyPairPEM()

	publicKey, err := GetPublicKey(publicPEM)

	if err != nil {
		t.Fatalf("GetPublicKey failed: %v", err)
	}

	// Verify it's a valid public key by using it to verify a signature
	privateKey, err := GetPrivateKey(privatePEM)
	if err != nil {
		t.Fatalf("GetPrivateKey failed: %v", err)
	}

	testData := []byte("test message")
	signature, err := GenerateAsymmetricSignature(privateKey, testData, crypto.SHA256)
	if err != nil {
		t.Fatalf("failed to generate signature: %v", err)
	}

	isValid, err := VerifyAsymmetricSignature(publicKey, testData, signature, crypto.SHA256)
	if err != nil {
		t.Fatalf("failed to verify signature: %v", err)
	}
	if !isValid {
		t.Fatal("signature verification failed")
	}
}

func TestPrivateAndPublicKey_KeyPairMatches(t *testing.T) {
	privatePEM, publicPEM := generateTestKeyPairPEM()

	privateKey, err := GetPrivateKey(privatePEM)
	if err != nil {
		t.Fatalf("GetPrivateKey failed: %v", err)
	}

	publicKey, err := GetPublicKey(publicPEM)
	if err != nil {
		t.Fatalf("GetPublicKey failed: %v", err)
	}

	// Verify the public key matches the private key
	if privateKey.PublicKey.N.Cmp(publicKey.N) != 0 {
		t.Fatal("public key does not match private key")
	}
	if privateKey.PublicKey.E != publicKey.E {
		t.Fatal("public key exponent does not match private key")
	}
}
