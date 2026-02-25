package security

import (
    "encoding/base64"
    "strings"
    "testing"
)

func TestGenerateRandomStaring(t *testing.T) {
    s1 := GenerateRandomStaring()
    if s1 == "" {
        t.Fatal("expected non-empty string")
    }
    // should contain only characters from allowed set and hex digits
    for _, ch := range s1 {
        if !strings.ContainsRune(Chars, ch) && !(ch >= '0' && ch <= '9') && !(ch >= 'a' && ch <= 'f') && !(ch >= 'A' && ch <= 'F') {
            t.Errorf("unexpected character %q in random string", ch)
        }
    }
    s2 := GenerateRandomStaring()
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

