package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"log"
	"os"
	"strings"
)

var (
	conn    *tls.Conn
	confVal []string
)

func init() {
	configFile := "./config.txt"
	configWords := []string{
		"manager_ip",
		"uptime_on",
		"lock_dir",
		"mon_interval",
	}
	searchKeywords(configFile, configWords)
	conn = initConn()
}

func main() {
	var mfaFlag string
	flag.StringVar(&mfaFlag, "mfa", "", "Give mfa candidate on command-line.")
	flag.Parse()

	if checkMfa(mfaFlag) {
		startSubroutine()
	}
}

func searchKeywords(filepath string, keywords []string) { //CONFVAL DIRECTORY
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, keyword := range keywords {
			if strings.Contains(line, "#") {
				break
			} else if strings.Contains(line, keyword) {
				result := strings.Split(line, "=")
				confVal = append(confVal, result[1])
			}
		}
	}
}

func initConn() *tls.Conn { //CONFVAL DIRECTORY
	config := &tls.Config{
		InsecureSkipVerify: true, // For self-signed certificates, in a production environment, set this to false and provide valid CA certificates.
	}

	conn, err := tls.Dial("tcp", confVal[0]+":9114", config)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	log.Println("Connected to", conn.RemoteAddr())
	return conn
}
