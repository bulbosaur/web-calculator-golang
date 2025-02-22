package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	orchestrator "github.com/bulbosaur/web-calculator-golang/internal/orchestrator/service"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
)

func RegHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
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

		_, err = orchestrator.Calc(request.Expression)
		if err != nil {
			exprRepo.UpdateStatus(id, models.StatusFailed)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: err.Error(),
			})
			return
		}

		result, err := orchestrator.Calc(request.Expression)
		if err != nil {
			exprRepo.UpdateStatus(id, models.StatusFailed)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: err.Error(),
			})
			return
		}

		exprRepo.UpdateStatus(id, models.StatusResolved)
		exprRepo.SetResult(id, result)

		response := models.RegisteredExpression{
			Id: id,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// // CalcHandler принимает json с выраженями и подсчитывает их значение при помощи Calc
// func CalcHandler(exprRepo *repository.ExpressionModel) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		request := new(models.Request)
// 		defer r.Body.Close()

// 		err := json.NewDecoder(r.Body).Decode(&request)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			json.NewEncoder(w).Encode(models.ErrorResponse{
// 				Error:        "Bad request",
// 				ErrorMessage: models.ErrorInvalidRequestBody.Error(),
// 			})
// 			return
// 		}

// 		id, err := exprRepo.Insert(request.Expression)
// 		if err != nil {
// 			log.Printf("something went wrong while creating a record in the database. %v", err)
// 		}
// 		log.Printf("Expression ID-%d has been registered", id)

// 		result, err := orchestrator.Calc(request.Expression)
// 		if err != nil {
// 			exprRepo.UpdateStatus(id, models.StatusFailed)
// 			w.WriteHeader(http.StatusUnprocessableEntity)
// 			json.NewEncoder(w).Encode(models.ErrorResponse{
// 				Error:        "Expression is not valid",
// 				ErrorMessage: err.Error(),
// 			})
// 			return
// 		}

// 		exprRepo.UpdateStatus(id, models.StatusResolved)

// 		response := models.Response{
// 			Result: result,
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(response)
// 	}
// }
