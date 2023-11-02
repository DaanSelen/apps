package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	rawInput, _ := ioutil.ReadFile("public.pem") // just pass the file name
	var publicPem string = string(rawInput)

	block, _ := pem.Decode([]byte(publicPem))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}

	fmt.Println("Enter the plaintext you wish to encrypt.")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	cipherText := RSA_OAEP_Encrypt(input, *publicKey)

	fmt.Println(cipherText)

	fmt.Scanln()
}

func RSA_OAEP_Encrypt(secretMessage string, key rsa.PublicKey) string {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, _ := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	return base64.StdEncoding.EncodeToString(ciphertext)
}
