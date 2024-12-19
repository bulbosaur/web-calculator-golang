package server

import (
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/config"
	"github.com/bulbosaur/web-calculator-golang/internal/handlers"
)

func RunServer() {
	config, err := config.GettingConfig(".../config/config.json")
	if err != nil {
		log.Fatal("Config error")
	}

	addr := config.Host + ":" + config.Port

	http.HandleFunc("/api/v1/calculate", handlers.CalcHandler)
	http.HandleFunc("/coffe", handlers.CoffeeHandler)

	log.Printf("Server starting on %s", addr)
	err = http.ListenAndServe(config.Port, nil)

	if err != nil {
		log.Fatal("Internal server error")
	}
}
