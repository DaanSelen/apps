package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func initHTTP() {
	fmt.Println("API-SERVER INITIALISING")
	CSMT := mux.NewRouter().StrictSlash(true)

	//Basic endpoints
	CSMT.HandleFunc("/", handleRootEndpoint).Methods("GET")
	CSMT.HandleFunc("/version", handleRootEndpoint).Methods("GET")

	//Control endpoints for the agent.
	CSMT.HandleFunc("/monitor/uptime", handleMonitorUptimeEndpoint).Methods("POST")
	CSMT.HandleFunc("/monitor/cpu", handleMonitorCPUAlertEndpoint).Methods("POST")
	CSMT.HandleFunc("/monitor/ram", handleMonitorRAMAlertEndpoint).Methods("POST")

	//Control endpoints (For looking up the values)
	CSMT.HandleFunc("/control/uptime", handleControlUptimeEndpoint).Methods("GET")
	CSMT.HandleFunc("/control/cpu", handleControlCPUEndpoint).Methods("GET")
	CSMT.HandleFunc("/control/ram", handleControlRAMEndpoint).Methods("GET")

	http.ListenAndServe((":2468"), CSMT)
}

// Basic endpoint
func handleRootEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	json.NewEncoder(w).Encode("Root directory endpoint hit! No options here.\nVERSION 0.09")
}

func handleMonitorUptimeEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var alert Alert
	json.NewDecoder(r.Body).Decode(&alert)
	if strings.ToLower(alert.Comp) == "uptime" {
		monitorUptimeAlert(alert)
	} else {
		w.WriteHeader(400)
	}
}

func handleMonitorCPUAlertEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var alert Alert
	json.NewDecoder(r.Body).Decode(&alert)
	if strings.ToLower(alert.Comp) == "cpu" {
		monitorCPUAlert(alert)
	} else {
		w.WriteHeader(400)
	}
}

func handleMonitorRAMAlertEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var alert Alert
	json.NewDecoder(r.Body).Decode(&alert)
	if strings.ToLower(alert.Comp) == "ram" {
		monitorRAMAlert(alert)
	} else {
		w.WriteHeader(400)
	}
}

func handleControlUptimeEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	alerts := checkComp("Uptime")
	if alerts == nil {
		json.NewEncoder(w).Encode("NOTHING RETURNED")
	} else {
		json.NewEncoder(w).Encode(alerts)
	}
}

func handleControlCPUEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	alerts := checkComp("CPU")
	if alerts == nil {
		json.NewEncoder(w).Encode("NOTHING RETURNED")
	} else {
		json.NewEncoder(w).Encode(alerts)
	}
}

func handleControlRAMEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	alerts := checkComp("RAM")
	if alerts == nil {
		json.NewEncoder(w).Encode("NOTHING RETURNED")
	} else {
		json.NewEncoder(w).Encode(alerts)
	}
}
