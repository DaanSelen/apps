package main

import (
	"log"
	"math/big"
	"net"
	"strings"
)

const (
	minInitVal = 200
	maxInitVal = 1000
)

func main() {
	const address string = "192.168.178.30:6800"

	s, _ := net.ResolveUDPAddr("udp4", address)
	connection, _ := net.ListenUDP("udp4", s)

	defer connection.Close()
	buffer := make([]byte, 1024)

	log.Print("SERVER IS READY FOR DHKEX, WAITING FOR MESSAGES")
	readIncomming(connection, buffer)
}

func readIncomming(conn *net.UDPConn, buffer []byte) {
	var active int
	initVals := [4]string{}

	for {
		n, addr, _ := conn.ReadFromUDP(buffer)
		if strings.TrimSpace(string(buffer[:n])) == "INITDIFFIE" {
			log.Print(strings.Repeat("-", 100))
			log.Println("ACTIVATED DHKEX")
			active = 1
		}

		if active > 0 && active < 5 {
			if active != 1 {
				initVals[active-2] = string(buffer[:n])
			}
			active++
			if active == 5 {
				initDiffie(conn, addr, initVals)
				active = 0
			}
		}

		if strings.TrimSpace(strings.ToUpper(string(buffer[:n]))) == "STOP" {
			log.Print("Exiting UDP server!")
			return
		} else if active == 0 && initVals[2] != string(buffer[:n]) { //Stop printing last DH value
			log.Print(string(buffer[:n]))
		}
	}
}

func initDiffie(conn *net.UDPConn, addr *net.UDPAddr, initVals [4]string) {
	b, _ := generateRandomPrime(big.NewInt(minInitVal), big.NewInt(maxInitVal))
	A, _ := new(big.Int).SetString(initVals[0], 10)
	g, _ := new(big.Int).SetString(initVals[1], 10)
	p, _ := new(big.Int).SetString(initVals[2], 10)
	log.Print("DIFFIE ACTIVATED RECEIVED INITIALIZATION VALUES")

	B := new(big.Int).Exp(g, b, p)

	data := []byte("Return:" + B.String())
	conn.WriteToUDP(data, addr)

	K := new(big.Int).Exp(A, b, p)
	log.Print("SHARED SECRET: ", K, " READY TO CREATE KEY")
	log.Println("VERBOSE VALUE INIT:", b, g, p)
	log.Println("VERBOSE VALUE WORK:", A, B)

	log.Println(strings.Repeat("-", 100))
}
