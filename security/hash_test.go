package security

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	s1 := GenerateRandomString()
	if s1 == "" {
		t.Fatal("expected non-empty string")
	}
	// should contain only characters from allowed set and hex digits
	for _, ch := range s1 {
		if !strings.ContainsRune(Chars, ch) && !(ch >= '0' && ch <= '9') && !(ch >= 'a' && ch <= 'f') && !(ch >= 'A' && ch <= 'F') {
			t.Errorf("unexpected character %q in random string", ch)
		}
	}
	s2 := GenerateRandomString()
	if s1 == s2 {
		t.Error("expected two generated strings to differ")
	}
}

func TestHmacShaKey(t *testing.T) {
	raw := []byte("some-secret")
	enc := base64.StdEncoding.EncodeToString(raw)

	got, err := HmacShaKey(enc)
	if err != nil {
		t.Fatalf("expected no error for valid base64, got %v", err)
	}
	if string(got) != string(raw) {
		t.Errorf("decoded bytes mismatch: got %q, want %q", got, raw)
	}

	_, err = HmacShaKey("not-base64-!!")
	if err == nil {
		t.Error("expected error for invalid base64 input")
	}
}

func TestHmacShaEncode(t *testing.T) {
	raw := "hello-world"
	enc := HmacShaEncode([]byte(raw))
	if enc == "" {
		t.Fatal("expected non-empty encoded string")
	}
	// verify it is base64
	decoded, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		t.Fatalf("encoded value is not valid base64: %v", err)
	}
	if string(decoded) != raw {
		t.Errorf("decoded mismatch: got %q, want %q", decoded, raw)
	}
}

func TestHmacShaKey_InvalidBase64(t *testing.T) {
	_, err := HmacShaKey("not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDigest_SHA256(t *testing.T) {
	data := []byte("test message")
	hash, err := Digest(crypto.SHA256, data)

	if err != nil {
		t.Fatalf("Digest failed: %v", err)
	}
	if len(hash) != 32 {
		t.Fatalf("SHA256 should produce 32 bytes, got %d", len(hash))
	}
}

func TestDigest_SHA512(t *testing.T) {
	data := []byte("test message")
	hash, err := Digest(crypto.SHA512, data)

	if err != nil {
		t.Fatalf("Digest failed: %v", err)
	}
	if len(hash) != 64 {
		t.Fatalf("SHA512 should produce 64 bytes, got %d", len(hash))
	}
}

func TestDigest_UnsupportedHashType(t *testing.T) {
	_, err := Digest(crypto.MD5, []byte("data"))
	if err == nil {
		t.Fatal("expected error for unsupported hash type")
	}
}

func TestHashBodyRequest_WithBytes(t *testing.T) {
	data := []byte("request body")
	hashed, err := HashBodyRequest(data, crypto.SHA256)

	if err != nil {
		t.Fatalf("HashBodyRequest failed: %v", err)
	}
	if hashed == "" {
		t.Fatal("expected non-empty hash")
	}

	// Verify it's valid hex
	_, err = hex.DecodeString(hashed)
	if err != nil {
		t.Fatalf("result is not valid hex: %v", err)
	}
}

func TestHashBodyRequest_WithStruct(t *testing.T) {
	type TestPayload struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}

	payload := TestPayload{Name: "test", ID: 123}
	hashed, err := HashBodyRequest(payload, crypto.SHA256)

	if err != nil {
		t.Fatalf("HashBodyRequest failed: %v", err)
	}
	if hashed == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestHashBodyRequest_Deterministic(t *testing.T) {
	data := []byte("consistent data")

	h1, err1 := HashBodyRequest(data, crypto.SHA256)
	h2, err2 := HashBodyRequest(data, crypto.SHA256)

	if err1 != nil || err2 != nil {
		t.Fatalf("HashBodyRequest failed: %v, %v", err1, err2)
	}
	if h1 != h2 {
		t.Fatalf("hashes should be identical for same input: %q != %q", h1, h2)
	}
}
