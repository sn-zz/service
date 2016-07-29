// sn - https://github.com/sn
package main

import (
    "crypto/rand"
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "golang.org/x/crypto/scrypt"
)

// GenerateUuid generates a universally unique identifier
func GenerateUuid() uuid {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        fmt.Println("Error: ", err)
        return ""
    }
    return uuid(fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
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
