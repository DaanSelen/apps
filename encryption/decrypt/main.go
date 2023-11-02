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
	rawInput, _ := ioutil.ReadFile("private.pem") // just pass the file name
	var privatePem string = string(rawInput)

	block, _ := pem.Decode([]byte(privatePem))
	if block == nil {
		panic("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}

	fmt.Println("Enter the text you wish to decrypt.")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	RSA_OAEP_Decrypt(input, *privateKey)

	fmt.Scanln()
}

func RSA_OAEP_Decrypt(cipherText string, privKey rsa.PrivateKey) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, _ := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	fmt.Println(string(plaintext))
}
