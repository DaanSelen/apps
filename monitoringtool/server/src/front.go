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

	go http.ListenAndServe((":9113"), NMTA)
}

func rootEnd(w http.ResponseWriter, r *http.Request) {
	log.Println("Root endpoint hit")
}

func accountMani(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch command {
		case "first option":
		case "second option":
		case "third option":
		}
	}
}