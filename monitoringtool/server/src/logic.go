package main

import (
	"fmt"
	"log"
)

func init() {
	initDB()
	initHTTP()
}

func main() {
	log.Println(infop, "Ready for incomming connections on port: "+listenPort)
	fmt.Scanln()
}
