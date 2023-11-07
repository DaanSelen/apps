package main

import (
	"fmt"
)

func init() {
	initDB()
	initHTTP()
}

func main() {
	fmt.Println("Ready")
	fmt.Scanln()
}


