package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestGenerateSymmetricSignature_Deterministic(t *testing.T) {
	secret := "s3cr3t"
	msg := []byte("the quick brown fox")
	sig1 := GenerateSymmetricSignature(secret, msg, crypto.SHA256)
	sig2 := GenerateSymmetricSignature(secret, msg, crypto.SHA256)

	if sig1 == "" {
		t.Fatal("expected non-empty signature")
	}
	if sig1 != sig2 {
		t.Fatalf("signatures differ for identical input: %q vs %q", sig1, sig2)
	}
}

func TestGenerateSymmetricSignature_VariesWithSecret(t *testing.T) {
	msg := []byte("data")
	s1 := GenerateSymmetricSignature("key-one", msg, crypto.SHA256)
	s2 := GenerateSymmetricSignature("key-two", msg, crypto.SHA256)

	if s1 == s2 {
		t.Fatalf("signatures should differ when secret changes: %q == %q", s1, s2)
	}
}

func TestGenerateSymmetricSignature_VariesWithMessage(t *testing.T) {
	secret := "topsecret"
	m1 := GenerateSymmetricSignature(secret, []byte("one"), crypto.SHA256)
	m2 := GenerateSymmetricSignature(secret, []byte("two"), crypto.SHA256)

	if m1 == m2 {
		t.Fatalf("signatures should differ when message changes: %q == %q", m1, m2)
	}
}

func TestGenerateSymmetricSignature_DifferentHashAlgorithms(t *testing.T) {
	secret := "hash-test"
	msg := []byte("payload")

	s256 := GenerateSymmetricSignature(secret, msg, crypto.SHA256)
	s512 := GenerateSymmetricSignature(secret, msg, crypto.SHA512)

	if s256 == s512 {
		t.Fatalf("signatures produced with different hash algorithms should differ: %q == %q", s256, s512)
	}
}

func TestGenerateAsymmetricSignature(t *testing.T) {
	// Generate a test RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	data := []byte("test message for asymmetric signature")
	sig, err := GenerateAsymmetricSignature(privateKey, data, crypto.SHA256)

	if err != nil {
		t.Fatalf("GenerateAsymmetricSignature failed: %v", err)
	}
	if sig == "" {
		t.Fatal("expected non-empty signature")
	}
}

func TestGenerateAsymmetricSignature_Deterministic(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	data := []byte("consistent data")

	sig1, err1 := GenerateAsymmetricSignature(privateKey, data, crypto.SHA256)
	sig2, err2 := GenerateAsymmetricSignature(privateKey, data, crypto.SHA256)

	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateAsymmetricSignature failed: %v, %v", err1, err2)
	}
	if sig1 != sig2 {
		t.Fatalf("signatures should be identical for same input: %q != %q", sig1, sig2)
	}
}

func TestVerifySymmetricSignature_Valid(t *testing.T) {
	secret := "shared-secret-key"
	msg := []byte("message to verify")

	// Generate a valid signature
	sig := GenerateSymmetricSignature(secret, msg, crypto.SHA256)

	// Verify it
	isValid, err := VerifySymmetricSignature(secret, msg, sig, crypto.SHA256)

	if err != nil {
		t.Fatalf("VerifySymmetricSignature failed: %v", err)
	}
	if !isValid {
		t.Fatal("expected signature to be valid")
	}
}

func TestVerifySymmetricSignature_InvalidSignature(t *testing.T) {
	secret := "shared-secret-key"
	msg := []byte("message to verify")

	// Verify with wrong signature
	isValid, err := VerifySymmetricSignature(secret, msg, "invalid-signature", crypto.SHA256)

	if err != nil {
		t.Fatalf("VerifySymmetricSignature failed: %v", err)
	}
	if isValid {
		t.Fatal("expected signature to be invalid")
	}
}

func TestVerifySymmetricSignature_WrongSecret(t *testing.T) {
	msg := []byte("message to verify")

	// Generate with one secret
	sig := GenerateSymmetricSignature("secret-one", msg, crypto.SHA256)

	// Verify with different secret
	isValid, err := VerifySymmetricSignature("secret-two", msg, sig, crypto.SHA256)

	if err != nil {
		t.Fatalf("VerifySymmetricSignature failed: %v", err)
	}
	if isValid {
		t.Fatal("expected signature to be invalid with wrong secret")
	}
}

func TestVerifySymmetricSignature_WrongMessage(t *testing.T) {
	secret := "shared-secret-key"

	// Generate with one message
	sig := GenerateSymmetricSignature(secret, []byte("original message"), crypto.SHA256)

	// Verify with different message
	isValid, err := VerifySymmetricSignature(secret, []byte("different message"), sig, crypto.SHA256)

	if err != nil {
		t.Fatalf("VerifySymmetricSignature failed: %v", err)
	}
	if isValid {
		t.Fatal("expected signature to be invalid with wrong message")
	}
}

func TestVerifySymmetricSignature_WrongHashAlgorithm(t *testing.T) {
	secret := "shared-secret-key"
	msg := []byte("message")

	// Generate with SHA256
	sig := GenerateSymmetricSignature(secret, msg, crypto.SHA256)

	// Verify with SHA512
	isValid, err := VerifySymmetricSignature(secret, msg, sig, crypto.SHA512)

	if err != nil {
		t.Fatalf("VerifySymmetricSignature failed: %v", err)
	}
	if isValid {
		t.Fatal("expected signature to be invalid with wrong hash algorithm")
	}
}

func TestVerifyAsymmetricSignature_Valid(t *testing.T) {
	// Generate a test RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	data := []byte("message to sign and verify")

	// Generate a valid signature
	sig, err := GenerateAsymmetricSignature(privateKey, data, crypto.SHA256)
	if err != nil {
		t.Fatalf("GenerateAsymmetricSignature failed: %v", err)
	}

	// Verify it with the public key
	isValid, err := VerifyAsymmetricSignature(publicKey, data, sig, crypto.SHA256)

	if err != nil {
		t.Fatalf("VerifyAsymmetricSignature failed: %v", err)
	}
	if !isValid {
		t.Fatal("expected signature to be valid")
	}
}

func TestVerifyAsymmetricSignature_InvalidSignature(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	data := []byte("message to verify")

	// Verify with invalid signature
	isValid, err := VerifyAsymmetricSignature(publicKey, data, "invalid-signature", crypto.SHA256)

	if err == nil {
		t.Fatal("expected error for invalid base64 signature")
	}
	if isValid {
		t.Fatal("expected signature to be invalid")
	}
}

func TestVerifyAsymmetricSignature_WrongMessage(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Sign one message
	sig, err := GenerateAsymmetricSignature(privateKey, []byte("original message"), crypto.SHA256)
	if err != nil {
		t.Fatalf("GenerateAsymmetricSignature failed: %v", err)
	}

	// Verify with different message
	isValid, err := VerifyAsymmetricSignature(publicKey, []byte("different message"), sig, crypto.SHA256)

	if err != nil {
		t.Fatalf("VerifyAsymmetricSignature failed: %v", err)
	}
	if isValid {
		t.Fatal("expected signature to be invalid for different message")
	}
}

func TestVerifyAsymmetricSignature_WrongHashAlgorithm(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	data := []byte("message")

	// Sign with SHA256
	sig, err := GenerateAsymmetricSignature(privateKey, data, crypto.SHA256)
	if err != nil {
		t.Fatalf("GenerateAsymmetricSignature failed: %v", err)
	}

	// Verify with SHA512
	isValid, err := VerifyAsymmetricSignature(publicKey, data, sig, crypto.SHA512)

	if err != nil {
		t.Fatalf("VerifyAsymmetricSignature failed: %v", err)
	}
	if isValid {
		t.Fatal("expected signature to be invalid with wrong hash algorithm")
	}
}

func TestVerifyAsymmetricSignature_WrongPublicKey(t *testing.T) {
	// Generate two different key pairs
	privateKey1, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	privateKey2, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	publicKey2 := &privateKey2.PublicKey

	data := []byte("message")

	// Sign with privateKey1
	sig, err := GenerateAsymmetricSignature(privateKey1, data, crypto.SHA256)
	if err != nil {
		t.Fatalf("GenerateAsymmetricSignature failed: %v", err)
	}

	// Verify with publicKey2 (wrong key)
	isValid, err := VerifyAsymmetricSignature(publicKey2, data, sig, crypto.SHA256)

	if err != nil {
		t.Fatalf("VerifyAsymmetricSignature failed: %v", err)
	}
	if isValid {
		t.Fatal("expected signature to be invalid with wrong public key")
	}
}
