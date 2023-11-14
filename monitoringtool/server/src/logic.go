package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	infop = "[Info]"
	warnp = "[Warn]"
	errop = "[Error]"

	JOINTOKEN_LEN = 100
)

func init() {
	log.Println(strings.Repeat("-", 100))
	initDB()
	initHTTP()
}

func main() {
	log.Println(infop, "Ready for incomming connections on port: "+listenPort+".")
	log.Println(strings.Repeat("-", 100))
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
	log.Println(infop, "Received request for account creation, user:", username+".")
	securedPassword, randomSalt := securePassword(password)
	joinToken := generateRandomString(JOINTOKEN_LEN)
	if status := insertAccount(username, securedPassword, randomSalt, joinToken); status {
		log.Println(infop, "Successfully created account for user:", username+".")
		return status
	} else {
		log.Println(warnp, "Failed to create account because of duplicate username.")
		return status
	}
}

func changeAccount(username, password, option string) bool {
	log.Println(infop, "Received request for account change, user:", username+".")
	if authenticateAccount(username, password) {
		log.Println(infop, "Passwords match, user:", username, "authenticated.")
		securedPassword, randomSalt := securePassword(option)
		alterAccount(username, securedPassword, randomSalt)
		log.Println(infop, "Account alteration succesful, altered password for user:", username+".")
		return true
	} else {
		log.Println(warnp, "Passwords do not match or user does not exist, access denied.")
		return false
	}
}

func removeAccount(username, password string) bool {
	log.Println(infop, "Received request for account deletion, user:", username+".")
	if authenticateAccount(username, password) {
		log.Println(infop, "Passwords match, user:", username, "authenticated.")
		dropAccount(username)
		log.Println(infop, "Account deletion succesful, removed user:", username+".")
		return true
	} else {
		log.Println(warnp, "Passwords do not match or user does not exist, access denied.")
		return false
	}
}

func getUserToken(username, password string) (bool, string) {
	log.Println(infop, "Received request for retrieval of accesstoken for user:", username+".")
	if authenticateAccount(username, password) {
		return true, retrieveUserToken(username)
	} else {
		return false, ""
	}
}

func registerAgent(agentManager, candidateAccessToken, agentHostname, agentOS, agentIP, agentSignDate string) bool {
	accessToken := retrieveUserToken(agentManager)
	if candidateAccessToken == accessToken {
		if insertAgent(agentManager, agentHostname, agentOS, agentIP, agentSignDate) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
