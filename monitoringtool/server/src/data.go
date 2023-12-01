package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

const (
	userTable = `
	CREATE TABLE IF NOT EXISTS users (
		id 			INTEGER PRIMARY KEY AUTOINCREMENT,
		user	 	TEXT,
		password 	TEXT,
		salt 		TEXT,
		accessToken TEXT
	);`
	agentTable = `
	CREATE TABLE IF NOT EXISTS agents (
		id 				INTEGER PRIMARY KEY AUTOINCREMENT,
		manager			TEXT,
		hostname 		TEXT,
		operatingSystem TEXT,
		ipAddress 		TEXT,
		signupDate		TEXT
	);`
	dataTable = `
	CREATE TABLE IF NOT EXISTS agents_data (
		id 			INTEGER PRIMARY KEY AUTOINCREMENT,
		hostname	TEXT,
		component 	TEXT,
		value 		TEXT
	);`

	userTableName   = "users"
	userColumnName  = "user"
	agentTableName  = "agents"
	agentColumnHost = "hostname"
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "NMASDB.db")
	if err != nil {
		log.Fatal(errop, err)
	}

	err = db.Ping() //Test the DB Connection to see if it works.
	if err != nil {
		log.Fatal(errop, "Failed to ping the database:", err)
	} else {
		db.Exec(userTable)  //Create first table containing account information, listed above.
		db.Exec(agentTable) //Create second table containing agent information, listed above.
		db.Exec(dataTable)  //Create third table containing agent data, listed above.

		log.Println(infop, "NMTAS SQLite3 Database, Ready for connections.")
	}
}

func checkDuplicate(keyword, columnname, table string) bool {
	var counter int
	query := "SELECT COUNT(*) FROM " + table + " WHERE " + columnname + " = ?;"
	err := db.QueryRow(query, keyword).Scan(&counter)
	if err != nil {
		log.Println("Error executing query:", err)
		return false
	}
	return counter > 0
}

func retrieveAmountOfAgents(remoteIP string) int {
	var counter int
	query := "SELECT COUNT(*) FROM " + agentTableName + " WHERE ipAddress = ?;"
	err := db.QueryRow(query, remoteIP).Scan(&counter)
	if err != nil {
		log.Println("Error executing query:", err)
	}
	return counter
}

func retrieveSalt(username string) (bool, string) {
	var randomSalt string
	query := "SELECT salt FROM users WHERE user = ?;"
	err := db.QueryRow(query, username).Scan(&randomSalt)
	if err != nil {
		log.Println("Error executing query:", err)
		return false, "Failed to find user: " + username
	}
	return true, randomSalt
}

func retrievePasswordhash(username string) string {
	var passwordHash string
	query := "SELECT password FROM users WHERE user = ?;"
	err := db.QueryRow(query, username).Scan(&passwordHash)
	if err != nil {
		log.Println("Error executing query:", err)
	}
	return passwordHash
}

func retrieveUserToken(username string) string {
	var accessToken string
	query := "SELECT accessToken FROM users WHERE user = ?;"
	err := db.QueryRow(query, username).Scan(&accessToken)
	if err != nil {
		log.Println("Error executing query:", err)
	}
	return accessToken
}

func insertAccount(username, securedPassword, randomSalt, joinToken string) bool {
	if !checkDuplicate(username, userColumnName, userTableName) {
		query := "INSERT INTO users (user, password, salt, accessToken) VALUES (?, ?, ?, ?);"
		_, err := db.Exec(query, username, securedPassword, randomSalt, joinToken)
		if err != nil {
			log.Println("Error executing query:", err)
			return false
		}
		return true
	}
	return false
}

func alterAccount(username, password, randomSalt string) {
	query := "UPDATE users SET password = ?, salt = ? WHERE user = ?;"
	_, err := db.Exec(query, password, randomSalt, username)
	if err != nil {
		log.Println("Error executing query:", err)
	}
}

func dropAccount(username string) {
	query := "DELETE FROM users WHERE user = ?;"
	_, err := db.Exec(query, username)
	if err != nil {
		log.Println("Error executing query:", err)
	}
}

func insertAgent(agentManager, agentHostname, agentOS, agentIP, signupDate string) bool {
	if !checkDuplicate(agentHostname, agentColumnHost, agentTableName) && agentManager != adminUsername {
		query := "INSERT INTO agents (manager, hostname, operatingSystem, ipAddress, signupDate) VALUES (?, ?, ?, ?, ?);"
		_, err := db.Exec(query, agentManager, agentHostname, agentOS, agentIP, signupDate)
		if err != nil {
			log.Println("Error executing query:", err)
			return false
		}
		return true
	}
	return false
}
