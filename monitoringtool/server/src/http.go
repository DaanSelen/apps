package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	listenPort = "9113"
	infop      = "[Info]"
	warnp      = "[Warn]"
	errop      = "[Error]"
)

type accountMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Option   string `json:"option"`
}

type infoMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

	log.Println(infop, "ROOT HIT") //Comment out later, for debugging purposes
	json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: "Nerthus Monitor Application Server REST-API. Version 0.01"})
}

func accountMani(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var requestBody accountMessage
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
		} else {
			switch command {
			case "create":
				status := createAccount(requestBody.Username, requestBody.Password)
				if status {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: ("Successfully created an account for user: " + requestBody.Username)}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusConflict)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusConflict, Message: "Creation failed, user: " + requestBody.Username + " exists."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case "change":
				log.Println(infop, "2 OPTON")
			case "remove":
				removeAccount(requestBody.Username, requestBody.Password)
			}
		}
	}
}
