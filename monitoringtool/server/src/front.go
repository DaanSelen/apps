package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initHTTP() {
	log.Println("NMTA Server Init")
	NMTA := mux.NewRouter().StrictSlash(true)

	//Basic endpoints
	NMTA.HandleFunc("/", rootEnd).Methods("GET")
	//Account endpoints
	NMTA.HandleFunc("/account/create", accountMani("create")).Methods("GET")
	NMTA.HandleFunc("/account/change", accountMani("change")).Methods("GET")
	NMTA.HandleFunc("/account/remove", accountMani("remove")).Methods("GET")

	go http.ListenAndServe((":9113"), NMTA)
}

func rootEnd(w http.ResponseWriter, r *http.Request) {
	log.Println("Root endpoint hit")
}

func accountMani(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch command {
		case "create":
			log.Println("1 OPTON")
		case "change":
			log.Println("2 OPTON")
		case "remove":
			log.Println("3 OPTON")
		}
	}
}
