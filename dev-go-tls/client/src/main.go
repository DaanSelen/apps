package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
)

func main() {
	serverAddr := "192.168.178.30"

	config := &tls.Config{
		InsecureSkipVerify: true, // For self-signed certificates, in a production environment, set this to false and provide valid CA certificates.
	}

	conn, err := tls.Dial("tcp", serverAddr+":9114", config)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	fmt.Println("Connected to", conn.RemoteAddr())

	// You can read and write data to the connection here.
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a message: ")
		scanner.Scan()
		message := scanner.Text()

		// Send the message to the server
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}

		// You can also implement logic to exit the loop based on some condition (e.g., type "exit")
		if message == "exit" {
			break
		}
	}
}
