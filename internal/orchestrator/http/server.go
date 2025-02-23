package orchestrator

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/repository"
	"github.com/spf13/viper"
)

// RunOrchestrator запускает сервер оркестратора
func RunOrchestrator(exprRepo, taskRepo *repository.ExpressionModel) {
	host := viper.GetString("server.ORC_HOST")
	port := viper.GetString("server.ORC_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	http.HandleFunc("POST /api/v1/calculate", RegHandler(exprRepo, taskRepo))
	http.HandleFunc("/internal/task", taskHandler(taskRepo))
	http.HandleFunc("/coffee", CoffeeHandler)

	log.Printf("Orchestrator starting on %s", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("Orchestrator server error:", err)
	}
}
