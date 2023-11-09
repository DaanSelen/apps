package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type infoMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	listenPort = "9113"
	infop      = "[Info]"
	warnp      = "[Warn]"
	errop      = "[Error]"
)

func initHTTP() {
	NMTA := mux.NewRouter().StrictSlash(true)

	//Basic endpoints
	NMTA.HandleFunc("/", rootEnd).Methods("GET")
	//Account endpoints
	NMTA.HandleFunc("/account/create", accountMani("create")).Methods("POST")
	NMTA.HandleFunc("/account/change", accountMani("change")).Methods("PATCH")
	NMTA.HandleFunc("/account/remove", accountMani("remove")).Methods("DELETE")
	//Agent register endpoint
	NMTA.HandleFunc("/agent/register", accountMani("create")).Methods("POST")
	NMTA.HandleFunc("/agent/deregister", accountMani("remove")).Methods("DELETE")

	go http.ListenAndServe((":" + listenPort), NMTA)
	log.Println(infop, "NMTAS HTTP REST-API, Ready for connections.")
}

func rootEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	returnMessage := infoMessage{Code: http.StatusOK, Message: "Nerthus Monitor Application Server REST-API. Version 0.01"}
	log.Println(infop, "ROOT HIT") //Comment out later, for debugging purposes
	json.NewEncoder(w).Encode(returnMessage)
}

func accountMani(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch command {
		case "create":
			log.Println(infop, "1 OPTON")
		case "change":
			log.Println(infop, "2 OPTON")
		case "remove":
			log.Println(infop, "3 OPTON")
		}
	}
}
