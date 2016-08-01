// sn - https://github.com/sn
package helpers

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"

	"github.com/sn/service/types"
)

// GenerateUUID generates a universally unique identifier
func GenerateUUID() types.UUID {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return types.UUID(fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
}

// GeneratePasswordHash generates a password hash using scrypt
func GeneratePasswordHash(password string) string {
	hash, err := scrypt.Key([]byte(password), []byte("!@)#(!@#"), 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// GenerateSha1Hash generates a sha1 hash
func GenerateSha1Hash(input string) string {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(input))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
