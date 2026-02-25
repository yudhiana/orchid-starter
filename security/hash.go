package security

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

const (
	Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GenerateRandomStaring() string {
	randomBytes := make([]byte, 16)
	for i := range randomBytes {
		randomBytes[i] = Chars[rand.Intn(len(Chars))]
	}
	timestampHex := fmt.Sprintf("%x", time.Now().UnixNano())
	minLength := len(randomBytes)
	if len(timestampHex) < minLength {
		minLength = len(timestampHex)
	}
	combined := make([]byte, 0, len(randomBytes)+len(timestampHex))
	for i := 0; i < minLength; i++ {
		combined = append(combined, randomBytes[i], timestampHex[i])
	}

	return string(combined)
}

func HmacShaKey(key string) (shaKey []byte, err error) {
	shaKey, errDecode := base64.StdEncoding.DecodeString(key)
	if errDecode != nil {
		err = fmt.Errorf("failed to decode key secret Error: %w", errDecode)
		return
	}
	return
}

func HmacShaEncode(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}
