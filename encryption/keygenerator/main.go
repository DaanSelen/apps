package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	checkError(err)
	publicKey := privateKey.PublicKey
	saveKeys(privateKey, &publicKey)
}

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func saveKeys(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) {
	// save PEM files
	privfile, err := os.Create("../decrypt/private.pem")
	checkError(err)
	defer privfile.Close()
	pubfile, err := os.Create("../encrypt/public.pem")
	checkError(err)
	defer pubfile.Close()

	var privPem = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey)}

	err = pem.Encode(privfile, privPem)
	checkError(err)

	var pubPem = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey)}

	err = pem.Encode(pubfile, pubPem)
	checkError(err)

}
