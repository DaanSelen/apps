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
	log.Println("Ready")
	fmt.Scanln()
}
