package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	createAgentTable string = `
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    email TEXT
	)`
)

func initDB() {
	db, err := sql.Open("sqlite3", "monitor.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Test the DB Connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	} else {
		db.Exec(createAgentTable) //Create first table containing agent information

		log.Println("FINISHING INSERTING FIRST TABLES")
	}
}
