package main

import (
	"crypto/tls"
	"fmt"
	"net"
)

const (
	tlsPort string = "6942"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Handle the connection here.
	// You can read and write data to the connection.

	fmt.Println("Connection established from", conn.RemoteAddr())

	// Example: Read data from the connection
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		fmt.Println("Received data:", string(buffer[:n]))
	}
}

func main() {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Println("Error loading server certificates:", err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", ":"+tlsPort, config)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on :" + tlsPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
