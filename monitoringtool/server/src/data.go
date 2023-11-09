package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const (
	userTable = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		password TEXT,
		salt TEXT
		);`
	agentTable = `
	CREATE TABLE IF NOT EXISTS agents (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	agentHostname TEXT,
    	signupDate TEXT,
		operatingSystem TEXT
	);`
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
		db.Exec(agentTable) //Create first table containing agent information, listed above.

		log.Println(infop, "NMTAS SQLite3 Database: Ready for connections.")
	}
}

func checkUserDuplicate(keyword string) bool {
	var counter int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE username = '" + keyword + "';").Scan(&counter)
	if counter == 0 {
		return false
	} else {
		return true
	}
}

func retrieveSalt(username string) (bool, string) {
	var randomSalt string
	db.QueryRow("SELECT salt FROM users WHERE username = '" + username + "';").Scan(&randomSalt)
	if len(randomSalt) == 0 {
		return false, "Failed to find user: " + username
	} else {
		return true, randomSalt
	}
}

func retrievePasswordhash(username string) string {
	var passwordHash string
	db.QueryRow("SELECT password FROM users WHERE username = '" + username + "';").Scan(&passwordHash)
	return passwordHash
}

func insertAccount(username, securedPassword, randomSalt string) bool {
	if !checkUserDuplicate(username) {
		stmnt, _ := db.Prepare("INSERT INTO users (username, password, salt) VALUES (?, ?, ?);")
		stmnt.Exec(username, securedPassword, randomSalt)
		return true
	} else {
		return false
	}
}

func alterAccount(username, password, randomSalt string) {
	log.Println(password, randomSalt)
	stmnt, _ := db.Prepare("UPDATE users SET password = ?, salt = ? WHERE username = ?;")
	defer stmnt.Close()
	stmnt.Exec(password, randomSalt, username)
}

func dropAccount(username string) {
	stmnt, _ := db.Prepare("DELETE FROM users WHERE username = ?;")
	defer stmnt.Close()
	stmnt.Exec(username)
}
