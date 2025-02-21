package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/handlers"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
	"github.com/spf13/viper"
)

// RunServer запускает сервер с заданными в config.yaml значениями
func RunServer(exprRepo *repository.ExpressionModel) {
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	addr := fmt.Sprintf("%s:%s", host, port)

	http.HandleFunc("POST /api/v1/calculate", handlers.CalcHandler(exprRepo))
	http.HandleFunc("GET /api/v1/expressions", handlers.ListHandler)
	http.HandleFunc("/coffee", handlers.CoffeeHandler)

	log.Printf("Server starting on %s", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("Internal server error")
	}
}
