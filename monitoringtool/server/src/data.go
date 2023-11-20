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
	CREATE TABLE IF NOT EXISTS agents (
    	id 			INTEGER PRIMARY KEY AUTOINCREMENT,
		hostname	TEXT,
    	component 	TEXT,
		value 		TEXT,
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
	db.QueryRow("SELECT COUNT(*) FROM " + table + " WHERE " + columnname + " = '" + keyword + "';").Scan(&counter)
	if counter == 0 {
		return false
	} else {
		return true
	}
}

func retrieveAmountOfAgents(remoteIP string) int {
	var counter int
	db.QueryRow("SELECT COUNT(*) FROM " + agentTableName + " WHERE ipAddress = '" + remoteIP + "';").Scan(&counter)
	return counter
}

func retrieveSalt(username string) (bool, string) {
	var randomSalt string
	db.QueryRow("SELECT salt FROM users WHERE user = '" + username + "';").Scan(&randomSalt)
	if len(randomSalt) == 0 {
		return false, "Failed to find user: " + username
	} else {
		return true, randomSalt
	}
}

func retrievePasswordhash(username string) string {
	var passwordHash string
	db.QueryRow("SELECT password FROM users WHERE user = '" + username + "';").Scan(&passwordHash)
	return passwordHash
}

func retrieveUserToken(username string) string {
	var accessToken string
	db.QueryRow("SELECT accesstoken FROM users WHERE user = '" + username + "';").Scan(&accessToken)
	return accessToken
}

func insertAccount(username, securedPassword, randomSalt, joinToken string) bool {
	if !checkDuplicate(username, userColumnName, userTableName) {
		stmnt, _ := db.Prepare("INSERT INTO users (user, password, salt, accessToken) VALUES (?, ?, ?, ?);")
		stmnt.Exec(username, securedPassword, randomSalt, joinToken)
		return true
	} else {
		return false
	}
}

func alterAccount(username, password, randomSalt string) {
	stmnt, _ := db.Prepare("UPDATE users SET password = ?, salt = ? WHERE user = ?;")
	defer stmnt.Close()
	stmnt.Exec(password, randomSalt, username)
}

func dropAccount(username string) {
	stmnt, _ := db.Prepare("DELETE FROM users WHERE user = ?;")
	defer stmnt.Close()
	stmnt.Exec(username)
}

func insertAgent(agentManager, agentHostname, agentOS, agentIP, signupDate string) bool {
	if !checkDuplicate(agentHostname, agentColumnHost, agentTableName) && agentManager != adminUsername { //Check if the requested registration is not a duplicate and to deny people trying to assign to the admin account.
		stmnt, _ := db.Prepare("INSERT INTO agents (manager, hostname, operatingSystem, ipAddress, signupDate) VALUES (?, ?, ?, ?, ?);")
		defer stmnt.Close()
		stmnt.Exec(agentManager, agentHostname, agentOS, agentIP, signupDate)
		return true
	} else {
		return false
	}
}
