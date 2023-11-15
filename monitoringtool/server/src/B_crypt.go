package main

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/sha3"
)

const (
	PW_SALT_LEN = 32
)

func generateRandomString(len int) string {
	randomSalt := make([]byte, len)
	rand.Read(randomSalt)
	return base64.StdEncoding.EncodeToString(randomSalt)
}

func securePassword(password string) (string, string) {
	randomSalt := generateRandomString(PW_SALT_LEN)
	securedPassword := generateHash(password + randomSalt)
	return securedPassword, randomSalt
}

func generateHash(candidate string) string {
	hash := sha3.New512()
	hash.Write([]byte(candidate))
	finishedHash := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(finishedHash)
}
