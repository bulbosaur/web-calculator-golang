package main

import (
	"database/sql"
	"log"

	config "github.com/bulbosaur/web-calculator-golang/config"
	server "github.com/bulbosaur/web-calculator-golang/internal/http"

	_ "modernc.org/sqlite"
)

func main() {
	config.Init()

	db, err := sql.Open("sqlite", "../db/calc.db")
	if err != nil {
		log.Fatalf("error when opening database: %d", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("error when connecting with database: %d", err)
	}

	log.Print("Successful connection to the database")
	defer db.Close()

	server.RunServer()
}
