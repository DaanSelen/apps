package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Alert struct {
	Hostname string `json:"hostname"`
	Comp     string `json:"comp"`
	IpAddr   string `json:"ipaddress"`
	Time     string `json:"time"`
}

const (
	timeFormat = "02-01-2006 15:04:05"

	initTag    = "[Init]"
	monitorTag = "[APPLICATION]"
	warningTag = "[WARNING]"
	errorTag   = "[Error]"
)

var (
	configKeywords = []string{"CPUENABLED", "RAMENABLED", "CPUTHRESHOLD", "RAMTHRESHOLD", "CPUALERTTIME", "RAMALERTTIME", "SERVERIP", "HTTPS"}

	cpuEnabled string
	ramEnabled string

	cpuThreshold float64
	cpuAlertTime int
	ramThreshold float64
	ramAlertTime int

	apiServerURL string

	ipAddr string
)

func main() {
	log.Println(initTag, "AGENT INITIALISING.")

	go checkRAMUsage()
	checkCPUUsage()
}

func sendReport(comp string) {
	hostname, _ := os.Hostname()
	t := time.Now().Format(timeFormat)
	alert := Alert{
		Hostname: hostname,
		Comp:     comp,
		IpAddr:   ipAddr,
		Time:     t,
	}
	body, _ := json.Marshal(alert)
	_, err := http.Post((apiServerURL + comp), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(errorTag, "FAILED TO SEND HTTP REQUEST\n", err)
	}
}
