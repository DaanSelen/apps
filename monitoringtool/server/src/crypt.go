package main

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/sha3"
)

const (
	PW_SALT_LEN = 32
)

func securePassword(password string) (string, string) {
	randomSalt := make([]byte, PW_SALT_LEN)
	rand.Read(randomSalt)
	securedPassword := generateHash(password + base64.StdEncoding.EncodeToString(randomSalt))
	return securedPassword, base64.StdEncoding.EncodeToString(randomSalt)
}

func generateHash(candidate string) string {
	hash := sha3.New512()
	hash.Write([]byte(candidate))
	finishedHash := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(finishedHash)
}
