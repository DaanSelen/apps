package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func init() { //Initialise the variables that will be necessary for operation.
	cpuEnabled = getInfoFromConfig(configKeywords[0])
	ramEnabled = getInfoFromConfig(configKeywords[1])
	cpuThreshold, _ = strconv.ParseFloat(getInfoFromConfig(configKeywords[2]), 64)
	ramThreshold, _ = strconv.ParseFloat(getInfoFromConfig(configKeywords[3]), 64)
	cpuAlertTime, _ = strconv.Atoi(getInfoFromConfig(configKeywords[4]))
	ramAlertTime, _ = strconv.Atoi(getInfoFromConfig(configKeywords[5]))
	if httpsEnabled := getInfoFromConfig(configKeywords[7]); strings.ToLower(httpsEnabled) == "true" {
		apiServerURL = "https://" + getInfoFromConfig(configKeywords[6]) + "/monitor/"
	} else {
		apiServerURL = "http://" + getInfoFromConfig(configKeywords[6]) + "/monitor/"
	}
	log.Println(initTag, "Full API-Server URL:", apiServerURL)

	ipAddr = getLocalIp()
}

func getInfoFromConfig(keyword string) string {
	var info string
	f, err := os.Open("config.txt")
	if err != nil {
		log.Println("Opening config.txt file, perhaps there is no config.txt?\n", err)
	}
	defer f.Close()

	for lineByLine := bufio.NewScanner(f); lineByLine.Scan(); {
		if !(strings.Contains(lineByLine.Text(), "#") || lineByLine.Text() == "") && strings.Contains(lineByLine.Text(), (keyword+" = ")) { //Skipping empty rows and commented rows and checking for the ' = ' in the config.
			info = strings.ReplaceAll(lineByLine.Text(), (keyword + " = "), "") //Getting the important information from the line with the keyword.
		}
	}
	return info
}

func getLocalIp() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	fullLocalAddr := conn.LocalAddr().(*net.UDPAddr)
	LocalAddr := strings.Split(fullLocalAddr.String(), ":")
	log.Println(initTag, "Local System IP ASSIGNED TO:", LocalAddr[0])
	return LocalAddr[0]
}
