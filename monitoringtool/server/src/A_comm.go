package main

import (
	"crypto/tls"
	"log"
	"net"
	"strings"
)

const (
	listenPortTls string = "9114"
)

func initTLS() {
	config := &tls.Config{
		Certificates: []tls.Certificate{loadTLSCertificate("tls.crt", "tls.key")},
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
			log.Println(errop, "Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
	if checkIfAllowedIP(remoteIP) {
		log.Println(infop, "Connection established from", remoteIP)

		for { // Example: Read data from the connection
			buffer := make([]byte, 1024)
			n, _ := conn.Read(buffer)
			data := string(buffer[:n])
			if len(data) < 1 || data == "exit" {
				conn.Close()
				break
			} else {
				switch {
				case strings.Contains(data, "uptime"):
					uptime := strings.TrimPrefix(data, "uptime:")
					log.Println(infop, "Uptime is:", uptime)
				case strings.Contains(data, "cpuutil"):
					cpuutil := strings.TrimPrefix(data, "cpuutil:")
					log.Println(infop, "Uptime is:", cpuutil)
				case strings.Contains(data, "ramutil"):
					ramutil := strings.TrimPrefix(data, "ramutil:")
					log.Println(infop, "Uptime is:", ramutil)
				default:
					log.Println(errop, "Received data was not defined:", data)
				}
			}
		}
	} else {
		log.Println(warnp, "Refused connection from:", remoteIP)
	}

}
