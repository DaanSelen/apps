package main

import (
	"fmt"
	"log"
	"strings"
)

func init() {
	initDB()
	initHTTP()
}

func main() {
	log.Println(infop, "Ready for incomming connections on port: "+listenPort+".\n"+strings.Repeat("-", 100))
	fmt.Scanln()
}

func authenticateAccount(username, password string) bool {
	status, randomSalt := retrieveSalt(username)
	if status {
		candidate := password + randomSalt
		candidateHash := generateHash(candidate)
		if candidateHash == retrievePasswordhash(username) {
			return true
		} else {
			return false
		}
	} else {
		log.Println(errop, randomSalt)
		return false
	}
}

func createAccount(username, password string) bool {
	log.Println(infop, "Received request for account creation, user:", username)
	securedPassword, randomSalt := securePassword(password)
	if status := insertAccount(username, securedPassword, randomSalt); status {
		log.Println(infop, "Successfully created account for user:", username+".")
		return status
	} else {
		log.Println(warnp, "Failed to create account because of duplicate username.")
		return status
	}
}

func removeAccount(username, password string) bool {
	log.Println(infop, "Received request for account deletion, user:", username+".")
	if authenticateAccount(username, password) {
		log.Println(infop, "Passwords match, user authenticated.")
		return true
	} else {
		return false
	}
}
