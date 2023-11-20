package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
)

const (
	infop = "[Info]"
	warnp = "[Warn]"
	errop = "[Error]"

	JOINTOKEN_LEN = 100

	adminUsername = "admin"
)

func init() {
	log.Println(strings.Repeat("-", 100))
	initDB()
	go createAdminAccount()
	go initTLS()
	go initHTTP()
}

func main() {
	fmt.Scanln()
}

func createAdminAccount() {
	adminToken := generateRandomString(JOINTOKEN_LEN)
	if status := insertAccount(adminUsername, "", "", adminToken); status {
		log.Println(infop, "Inserted Admin account. Access token:")
		log.Println(infop, adminToken)
	} else {
		log.Println(infop, "Admin account already exists, not redoing.")
	}
}

func loadTLSCertificate(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	return cert
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

func checkIfAgentExists(remoteIP string) bool {
	return retrieveAmountOfAgents(remoteIP) > 0
}

func createAccount(username, password, option string) bool {
	log.Println(infop, "Received request for account creation, user:", username+".")
	securedPassword, randomSalt := securePassword(password)
	joinToken := generateRandomString(JOINTOKEN_LEN)
	if option == retrieveUserToken(adminUsername) {
		if status := insertAccount(username, securedPassword, randomSalt, joinToken); status {
			log.Println(infop, "Successfully created account for user:", username+".")
			return status
		} else {
			log.Println(warnp, "Failed to create account because of duplicate username.")
			return status
		}
	} else {
		return false
	}
}

func changeAccount(username, password, option string) bool {
	log.Println(infop, "Received request for account change, user:", username+".")
	if authenticateAccount(username, password) {
		log.Println(infop, "Passwords match, user:", username, "authenticated.")
		securedPassword, randomSalt := securePassword(option)
		alterAccount(username, securedPassword, randomSalt)
		log.Println(infop, "Account alteration successful, altered password for user:", username+".")
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
		log.Println(infop, "Account deletion successful, removed user:", username+".")
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

func checkIfAllowedIP(remoteIP string) bool {
	if checkIfAgentExists(remoteIP) {
		log.Println(infop, "Agent exists in database and is allowed, accepting connection.")
		return true
	} else {
		return false
	}
}
