package main

import (
	"log"

	config "github.com/bulbosaur/web-calculator-golang/config"
	server "github.com/bulbosaur/web-calculator-golang/internal/http"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"

	_ "modernc.org/sqlite"
)

func main() {
	config.Init()

	db, err := repository.InitDB("../db/calc.db")
	if err != nil {
		log.Fatal(err)
	}

	exprRepo := repository.NewExpressionModel(db)

	defer db.Close()

	server.RunServer(exprRepo)
}
