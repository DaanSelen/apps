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
	log.Println("GOT VALUES:", username, password)
	securedPassword, randomSalt := securePassword(password)
	return insertAccount(username, securedPassword, randomSalt)
}
