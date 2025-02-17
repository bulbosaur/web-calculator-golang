package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/pkg/calc"
)

// CalcHandler принимает json с выраженями и подсчитывает их значение при помощи Calc
func CalcHandler(w http.ResponseWriter, r *http.Request) {
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

	result, err := calc.Calc(request.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:        "Expression is not valid",
			ErrorMessage: err.Error(),
		})
		return
	}

	response := models.Response{
		Result: result,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
