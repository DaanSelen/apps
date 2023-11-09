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
		)`
	agentTable = `
	CREATE TABLE IF NOT EXISTS agents (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	agentHostname TEXT,
    	signupDate TEXT,
		operatingSystem TEXT
	)`
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

func insertAccount(username, securedPassword, randomSalt string) string {
	if !checkUserDuplicate(username) {
		statement, _ := db.Prepare("INSERT INTO users (username, password, salt) VALUES (?, ?, ?)")
		statement.Exec(username, securedPassword, randomSalt)
		return "SUCCESS"
	} else {
		return "FAILED, USERNAME EXISTS"
	}
}
