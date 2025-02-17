package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/config"
	"github.com/bulbosaur/web-calculator-golang/internal/handlers"
)

func RunServer() {
	config := config.GettingConfig()
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	http.HandleFunc("POST /api/v1/calculate", handlers.CalcHandler)
	http.HandleFunc("/coffee", handlers.CoffeeHandler)

	log.Printf("Server starting on %s", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("Internal server error")
	}
}
