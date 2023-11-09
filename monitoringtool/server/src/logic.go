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
	log.Println(infop, "Ready for incomming connections on port: "+listenPort+"\n"+strings.Repeat("-", 100))
	fmt.Scanln()
}

func createAccount(username, password string) string {
	log.Println(infop, "Received request for account creation, user:", username)
	securedPassword, randomSalt := securePassword(password)
	if status := insertAccount(username, securedPassword, randomSalt); status == "SUCCESS" {
		log.Println(infop, "Successfully created account for user:", username)
		return status
	} else {
		log.Println(warnp, "Failed to create account because of duplicate username.")
		return status
	}
}

/*func authenticateAccount() {

}*/
