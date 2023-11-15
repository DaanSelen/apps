package main

import (
	"crypto/tls"
	"log"
	"net"
)

const (
	listenPortTls string = "9114"
)

func initTLS() {

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println("Error loading server certificates:", err)
		return
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", ":"+listenPortTls, config)
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	log.Println(infop, "NMTAS TLS-Server, Ready for connections on port:", listenPortTls)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Connection established from", conn.RemoteAddr())

	for { // Example: Read data from the connection
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}

		log.Println("Received data:", string(buffer[:n]))
	}
}
