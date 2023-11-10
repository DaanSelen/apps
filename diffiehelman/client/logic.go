package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

var (
	B *big.Int
)

func main() {
	var CONNECT = "192.168.178.30:6800"

	s, _ := net.ResolveUDPAddr("udp4", CONNECT)
	conn, _ := net.DialUDP("udp4", nil, s)

	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	log.Println("READY FOR DHKEX, WAITING FOR MESSAGES")

	go readRespone(conn)

	for {

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		if strings.TrimSpace(strings.ToUpper(string(data))) == "START" {
			initDiffie(conn)
		} else {
			conn.Write(data)
			if strings.TrimSpace(strings.ToUpper(string(data))) == "STOP" {
				fmt.Println("Exiting UDP client!")
				return
			}
		}
	}
}

func readRespone(conn *net.UDPConn) {
	for {
		buffer := make([]byte, 1024)
		n, _, _ := conn.ReadFromUDP(buffer)
		if strings.Contains(string(buffer[:n]), "Return") {
			B, _ = new(big.Int).SetString(string((buffer[:n])[7:]), 10)
			log.Println("GOT PEER RETURN VALUE:", B)
		}
	}
}

func initDiffie(conn *net.UDPConn) {
	a, _ := generateRandomPrime(big.NewInt(50), big.NewInt(200))
	g, _ := generateRandomPrime(big.NewInt(50), big.NewInt(200))
	p, _ := generateRandomPrime(big.NewInt(100), big.NewInt(200))

	A := new(big.Int).Exp(g, a, p)

	log.Println("ACTIVATING DHKEX")

	conn.Write([]byte("INITDIFFIE"))
	vals := []*big.Int{A, g, p}
	for _, value := range vals {
		data := []byte(value.String())
		conn.Write(data)
	}

	time.Sleep(100 * time.Nanosecond)
	if B != nil {
		K := new(big.Int).Exp(B, a, p)
		log.Println("SHARED SECRET:", K, "READY TO CREATE KEY")
		log.Println("VERBOSE VALUE INIT:", a, g, p)
		log.Println("VERBOSE VALUE WORK:", A, B)
	} else {
		log.Println("B NOT RECEIVED.")
	}
}
