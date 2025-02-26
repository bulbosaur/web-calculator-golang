package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
)

func listHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expr, err := exprRepo.GetExpression(1)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.Expression{
			ID:     expr.ID,
			Status: expr.Status,
			Result: expr.Result,
		})
	}
}
