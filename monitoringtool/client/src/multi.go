package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	lockpath string = "/tmp"
)

func checkMfa(mfaCand string) bool {
	lockPathVerify()

	log.Println("Checking if MFA has been verified.")
	if _, err := os.Stat(lockpath + "/.nmac_mfa.lock"); err == nil {
		log.Println("MFA Lock found, not redoing.")
		return true
	} else {
		if mfaVerify(mfaCand) {
			os.Create(lockpath + "/.nmac_mfa.lock")
			return true
		} else {
			return false
		}
	}
}

func lockPathVerify() { //CONFVAL DIRECTORY
	if confVal[2] == "" {
		log.Println("Default lockpath chosen")
	} else if strings.HasPrefix(confVal[2], "/") {
		log.Println("Custom lockpath chosen")
		lockpath = confVal[1]
	}
}

func mfaVerify(mfaCand string) bool {
	var message string
	if mfaCand == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("mfa:")
		scanner.Scan()
		message = scanner.Text()
	} else {
		message = mfaCand
	}
	conn.Write([]byte("mfa:" + message))

	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	data := string(buffer[:n])
	if data == "SUCCESS" {
		log.Println("SUCCESSFUL MFA VERIFICATION")
		return true
	} else {
		log.Println("UNSUCCESSFUL MFA VERIFICATION")
		return false
	}
}
