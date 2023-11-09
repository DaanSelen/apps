package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	accountTable string = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		password TEXT,
		 TEXT,
		)`
	agentTable string = `
	CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	agentHostname TEXT,
    	signupDate TEXT,
		operatingSystem TEXT,
	)`
)

func initDB() {
	db, err := sql.Open("sqlite3", "NMASDB.db")
	if err != nil {
		log.Fatal(errop, err)
	}
	defer db.Close()

	err = db.Ping() //Test the DB Connection to see if it works.
	if err != nil {
		log.Fatal(errop, "Failed to ping the database:", err)
	} else {
		db.Exec(accountTable) //Create first table containing account information, listed above.
		db.Exec(agentTable)   //Create first table containing agent information, listed above.

		log.Println(infop, "NMTAS SQLite3 Database: Ready for connections.")
	}
}
