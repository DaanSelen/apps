package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	listenAddrHttp string = "0.0.0.0:9113"
	restApiTLSCert string = "./certs/restapi.crt"
	restApiTLSKey  string = "./certs/restapi.key"
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
	httpServer := &http.Server{
		Addr:    listenAddrHttp, // Specify the desired HTTPS port
		Handler: NMTA,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{ // Load the certificate and private key
				loadTLSCertificate(restApiTLSCert, restApiTLSKey),
			},
		},
	}

	//Basic endpoints
	NMTA.HandleFunc("/", rootEnd).Methods("GET")
	//Account endpoints
	NMTA.HandleFunc("/account/create", accountMani(1)).Methods("POST")
	NMTA.HandleFunc("/account/change", accountMani(2)).Methods("PATCH")
	NMTA.HandleFunc("/account/remove", accountMani(3)).Methods("DELETE")
	NMTA.HandleFunc("/account/accesstoken", accountMani(4)).Methods("GET")
	NMTA.HandleFunc("/account/mfa", accountMani(5)).Methods("GET")
	//Agent register endpoint
	NMTA.HandleFunc("/agent/register", agentMani("register")).Methods("POST")
	NMTA.HandleFunc("/agent/deregister", agentMani("deregister")).Methods("DELETE")

	go httpServer.ListenAndServeTLS("", "")
	log.Println(infop, "NMTAS HTTPS REST-API, Ready for connections on port:", listenAddrHttp)
}

func rootEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Println(infop, "ROOT HIT") //Comment out later, for debugging purposes
	json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: "Nerthus Monitor Application Server REST-API. Version 0.3.0"})
}

func accountMani(command int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var requestBody accountMessage
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		} else if requestBody.Username != adminUsername {
			switch command {
			case 1: //Create a new account and store the given password (or password hash) securely with an added salt using SHA3-512.
				if createAccount(requestBody.Username, requestBody.Password, requestBody.Option) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: ("Successfully created an account for user: " + requestBody.Username + ".")}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusConflict)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusConflict, Message: "Creation failed, user: " + requestBody.Username + " exists. Or the create token is incorrect."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case 2:
				if changeAccount(requestBody.Username, requestBody.Password, requestBody.Option) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: ("Successfully changed password for user: " + requestBody.Username + ".")}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusUnauthorized, Message: "Account change failed, user does not exist or credentials are incorrect."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case 3: //Check if the entered credentials are (when rehashed) equal to the stored credentials, if correct initiate account deletion.
				if removeAccount(requestBody.Username, requestBody.Password) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: ("Successfully removed account: " + requestBody.Username + ".")}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusUnauthorized, Message: "Deletion failed, user does not exist or credentials are incorrect."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case 4: //Check if the entered credentials are (when rehashed) equal to the stored credentials, if correct give the access token to register agents
				if status, userToken := getUserToken(requestBody.Username, requestBody.Password); status {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: userToken}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusUnauthorized, Message: "Retrieval of access code failed, user does not exist or credentials are incorrect."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case 5:
				if status, totpToken := getMFAToken(requestBody.Username, requestBody.Password); status {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: totpToken}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusUnauthorized, Message: "Retrieval of multifactor code failed, user does not exist or credentials are incorrect."}) //Using the predefined struct above we respond in JSON to the request.
				}
			}
		} else {
			log.Println(warnp, "Received calls for actions regarding the admin user, dropping.")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(infoMessage{Code: http.StatusUnauthorized, Message: "Dropping requests for the admin user."})
		}
	}
}

func agentMani(command string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var requestBody agentMessage
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		} else {
			switch command {
			case "register":
				if registerAgent(requestBody.AgentManager, requestBody.AccessToken, requestBody.AgentHostname, requestBody.AgentOS, requestBody.AgentIP, requestBody.AgentDate) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: "Succesfully registered agent to manager."}) //Using the predefined struct above we respond in JSON to the request.
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(infoMessage{Code: http.StatusOK, Message: "Access token was incorrect, manager account does not exist or agent is a duplicate."}) //Using the predefined struct above we respond in JSON to the request.
				}
			case "deregister":
				//deregisterAgent(requestBody.AgentHostname, requestBody.AgentOS)
				log.Println("Deregister")
			}
		}
	}
}
