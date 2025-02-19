package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
	"github.com/bulbosaur/web-calculator-golang/pkg/calc"
)

func CalcHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := new(models.Request)
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Bad request",
				ErrorMessage: models.ErrorInvalidRequestBody.Error(),
			})
			return
		}

		id, err := exprRepo.Insert(request.Expression)
		if err != nil {
			log.Printf("something went wrong while creating a record in the database. %v", err)
		}
		log.Printf("Expression ID-%d has been registered", id)

		result, err := calc.Calc(request.Expression)
		if err != nil {
			exprRepo.UpdateStatus(id, models.StatusFailed)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: err.Error(),
			})
			return
		} else {
			exprRepo.UpdateStatus(id, models.StatusResolved)
		}

		response := models.Response{
			Result: result,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
