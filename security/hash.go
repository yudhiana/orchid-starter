package security

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
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

func HmacShaEncode(raw []byte) string {
	return base64.StdEncoding.EncodeToString(raw)
}

func Digest(hashType crypto.Hash, signer []byte) ([]byte, error) {
	switch hashType {
	case crypto.SHA256:
		hash := crypto.SHA256.New()
		hash.Write(signer)
		return hash.Sum(nil), nil
	case crypto.SHA512:
		hash := crypto.SHA512.New()
		hash.Write(signer)
		return hash.Sum(nil), nil
	default:
		return nil, fmt.Errorf("unsupported hash type: %v", hashType)
	}
}

func HashBodyRequest(bodyReq any, hashType crypto.Hash) (bodyHashed string, err error) {
	var body []byte
	switch input := bodyReq.(type) {
	case []byte:
		body = input
	default:
		dataMarshal, errMarshal := json.Marshal(bodyReq)
		if errMarshal != nil {
			err = fmt.Errorf("failed marshal body request Error: %w", errMarshal)
			return
		}
		body = dataMarshal
	}

	hashed, errHash := Digest(hashType, body)
	if errHash != nil {
		err = fmt.Errorf("failed to hash body request Error: %w", errHash)
		return
	}
	bodyHashed = strings.ToLower(hex.EncodeToString(hashed))
	return
}
