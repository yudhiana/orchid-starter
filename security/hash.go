package security

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/mataharibiz/sange/v2"
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
		errMsg := fmt.Errorf("failed to decode secret. Error: %w", errDecode)
		return nil, sange.NewSetError(sange.OperationFailed, errMsg, errMsg.Error(), "security", "HmacShaKey")
	}

	return
}
