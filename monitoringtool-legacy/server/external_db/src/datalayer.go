package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ipaddr      = "localhost:3306"
	auth        = "celdserv:Ds190703"
	workingbase = "CSMT"
)

var (
	entry *sql.DB
	err   error
)

func initDBConnection() {
	fmt.Println("ATTEMPTING DATABASE CONNECTION")
	entry, err = sql.Open("mysql", auth+"@tcp("+ipaddr+")/"+workingbase)
	if err != nil {
		log.Println("[ERROR] CONNECTING TO DATABASE FAILED")
	}
	entry.Query("CREATE TABLE IF NOT EXISTS cpuentry (id INTEGER PRIMARY KEY AUTO_INCREMENT, hostname TEXT NOT NULL, ipaddress TEXT NOT NULL, comp TEXT NOT NULL, time TEXT NOT NULL)")
	entry.Query("CREATE TABLE IF NOT EXISTS ramentry (id INTEGER PRIMARY KEY AUTO_INCREMENT, hostname TEXT NOT NULL, ipaddress TEXT NOT NULL, comp TEXT NOT NULL, time TEXT NOT NULL)")

	checkConn()
}

func checkConn() {
	var currentEntries int

	data, err := entry.Query("select count(id) from cpuentry")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	data.Next()
	data.Scan(&currentEntries)
	fmt.Println("Current entries:", currentEntries)
}

func insertEntry(alert Alert) {
	entry.Query("INSERT INTO " + alert.Comp + "entry(hostname, ipaddress, comp, time) VALUES('" + alert.Hostname + "', '" + alert.IpAddr + "', '" + alert.Comp + "', '" + alert.Time + "')")
}

func dataCheck(component string) []Alert {
	var alerts []Alert
	rows, _ := entry.Query("SELECT * FROM " + component + "entry")
	defer rows.Close()
	for rows.Next() {
		var alert Alert
		rows.Scan(&alert.ID, &alert.Hostname, &alert.IpAddr, &alert.Comp, &alert.Time)
		alerts = append(alerts, alert)
	}

	return alerts
}
