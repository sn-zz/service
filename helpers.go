// sn - https://github.com/sn
package main

import (
    "crypto/rand"
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "golang.org/x/crypto/scrypt"
)

func GenerateUuid() uuid {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        fmt.Println("Error: ", err)
        return ""
    }
    return uuid(fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
}

func GeneratePasswordHash(password string) string {
    hash, _ := scrypt.Key([]byte(password), []byte("!@)#(!@#"), 16384, 8, 1, 32)
    return string(hash)
}

func GenerateSha1Hash(input string) string {
    hasher := sha1.New()
    hasher.Write([]byte(input))
    return hex.EncodeToString(hasher.Sum(nil))
}
