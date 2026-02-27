package security

import (
	"crypto"
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
