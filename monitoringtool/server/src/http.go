package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	listenPort = "9113"
)

type infoMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type accountMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Option   string `json:"option"`
}

type agentMessage struct {
	AgentManager  string `json:"agentmanager"`
	AccessToken   string `json:"accesstoken"`
	AgentHostname string `json:"agenthostname"`
	AgentOS       string `json:"agentos"`
	AgentIP       string `json:"agentip"`
	AgentDate     string `json:"agentdate"`
}

func initHTTP() {
	NMTA := mux.NewRouter().StrictSlash(true)

	//Basic endpoints
	NMTA.HandleFunc("/", rootEnd).Methods("GET")
	//Account endpoints
	NMTA.HandleFunc("/account/create", accountMani("create")).Methods("POST")
	NMTA.HandleFunc("/account/change", accountMani("change")).Methods("PATCH")
	NMTA.HandleFunc("/account/remove", accountMani("remove")).Methods("DELETE")
	NMTA.HandleFunc("/account/accesstoken", accountMani("accesstoken")).Methods("GET")
	//Agent register endpoint
	NMTA.HandleFunc("/agent/register", agentMani("register")).Methods("POST")
	NMTA.HandleFunc("/agent/deregister", agentMani("deregister")).Methods("DELETE")

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
			case "create": //Create a new account and store the given password (or password hash) securely with an added salt using SHA3-512.
				status := createAccount(requestBody.Username, requestBody.Password)
				if status {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: ("Successfully created an account for user: " + requestBody.Username + ".")}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusConflict)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusConflict, Message: "Creation failed, user: " + requestBody.Username + " exists."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case "change":
				if changeAccount(requestBody.Username, requestBody.Password, requestBody.Option) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: ("Successfully changed password for user: " + requestBody.Username + ".")}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusUnauthorized, Message: "Account change failed, user does not exist or credentials are incorrect."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case "remove": //Check if the entered credentials are (when rehashed) equal to the stored credentials, if correct initiate account deletion.
				if removeAccount(requestBody.Username, requestBody.Password) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: ("Successfully removed account: " + requestBody.Username + ".")}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusUnauthorized, Message: "Deletion failed, user does not exist or credentials are incorrect."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case "accesstoken":
				if status, userToken := getUserToken(requestBody.Username, requestBody.Password); status {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: userToken}) //Using the predefined struct above we respond in JSON to the request.
				}
			}
		}
	}
}

func agentMani(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var requestBody agentMessage
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
		} else {
			switch command {
			case "register":
				if registerAgent(requestBody.AgentManager, requestBody.AccessToken, requestBody.AgentHostname, requestBody.AgentOS, requestBody.AgentIP, requestBody.AgentDate) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: "Succesfully registered agent to manager."}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					log.Println("Unauthorized")
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: "Access token was incorrect or manager account does not exist."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case "deregister":
				//deregisterAgent(requestBody.AgentHostname, requestBody.AgentOS)
				log.Println("Deregister")
			}
		}
	}
}
