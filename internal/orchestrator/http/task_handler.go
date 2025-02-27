package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
)

func taskHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var task *models.Task

			task, id, err := exprRepo.GetTask()
			if err != nil {
				log.Println("Failed to get task:", err)
				http.Error(w, "Failed to get task", http.StatusInternalServerError)
				return
			}

			if task == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"task": task})

			exprRepo.UpdateTaskStatus(id, models.StatusCalculate)

		case http.MethodPost:
			var req struct {
				Id     int     `json:"id"`
				Result float64 `json:"result"`
				Error  error   `json:"error"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
				return
			}

			err = exprRepo.UpdateTaskResult(req.Id, req.Result, req.Error)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
