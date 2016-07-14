// sn - https://github.com/sn
package main

import (
    "crypto/rand"
    "fmt"
    "golang.org/x/crypto/scrypt"
)

func GenerateUuid() string {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        fmt.Println("Error: ", err)
        return ""
    }
    return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func GenerateHash(password string) string {
    hash, _ := scrypt.Key([]byte(password), []byte("!@)#(!@#"), 16384, 8, 1, 32)
    return string(hash)
}
