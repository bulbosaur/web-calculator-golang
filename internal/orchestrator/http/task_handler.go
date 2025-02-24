package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
)

func taskHandler(taskRepo *repository.ExpressionModel) http.HandlerFunc {
	var result struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var task *models.Task

			task, err := taskRepo.GetTask()
			if err != nil {
				log.Println("Failed to get task:", err)
				http.Error(w, "Failed to get task", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"task": task})

		case http.MethodPost:
			var req struct {
				Id     int     `json:"id"`
				Result float64 `json:"result"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
				return
			}

			// taskRepo.UpdateStatus(req.Id, models.StatusResolved)
			if err := taskRepo.UpdateTaskResult(result.ID, result.Result); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
