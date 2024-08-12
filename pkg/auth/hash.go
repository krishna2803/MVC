package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func CreateHash(password string) (string, error) {
	memory := uint32(64 * 1024)
	time := uint32(1)
	parallelism := uint8(4)
	keyLength := uint32(32)

	// salt changes every time
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}

	// argon2 ID hash, invulnerable to side-channel attacks
	hash := argon2.IDKey([]byte(password), []byte(salt), time, memory, parallelism, keyLength)

	// salt<64>.hash<64>
	return base64.RawStdEncoding.EncodeToString(salt) + "." + base64.RawStdEncoding.EncodeToString(hash), nil
}

func VerifyHash(password string, hash string) (bool, error) {
	hashParts := strings.Split(hash, ".")
	if len(hashParts) != 2 {
		return false, fmt.Errorf("invalid hash")
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[0])
	if err != nil {
		return false, err
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[1])
	if err != nil {
		return false, err
	}

	comparisonHash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)

	return bytes.Equal(decodedHash, comparisonHash), nil
}
