package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
)

func TestCalcHandler(t *testing.T) {
	cases := []struct {
		requestBody models.Request
		wantStatus  int
		wantOutput  *models.Response
		wantError   *models.ErrorResponse
	}{
		{
			requestBody: models.Request{
				Expression: "2+2",
			},
			wantStatus: http.StatusOK,
			wantOutput: &models.Response{
				Result: 4,
			},
			wantError: nil,
		},
		{
			requestBody: models.Request{
				Expression: "(2+2)*2",
			},
			wantStatus: http.StatusOK,
			wantOutput: &models.Response{
				Result: 8,
			},
			wantError: nil,
		},
		{
			requestBody: models.Request{
				Expression: "2++2",
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantOutput: nil,
			wantError: &models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: models.ErrorInvalidInput.Error(),
			},
		},
		{
			requestBody: models.Request{
				Expression: "1/0",
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantOutput: nil,
			wantError: &models.ErrorResponse{
				Error:        "Expression is not valid",
				ErrorMessage: models.ErrorDivisionByZero.Error(),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.requestBody.Expression, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			rr := httptest.NewRecorder()

			db, err := repository.InitDB("../db/calc.db")
			if err != nil {
				t.Fatalf("failed to init DB: %v", err)
			}

			handler := http.HandlerFunc(CalcHandler(repository.NewExpressionModel(db)))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("%v: CalcHandler returned %v, but want %v", tc.requestBody.Expression, status, tc.wantStatus)
			}

			if tc.wantOutput != nil {
				var response models.Response
				err = json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				if response.Result != tc.wantOutput.Result {
					t.Errorf("%v: CalcHandler returned %v, but want %v", tc.requestBody.Expression, response.Result, tc.wantOutput.Result)
				}
			}

			if tc.wantError != nil {
				var errorResponse models.ErrorResponse
				err = json.NewDecoder(rr.Body).Decode(&errorResponse)
				if err != nil {
					t.Fatalf("failed to decode error response body: %v", err)
				}

				if errorResponse.Error != tc.wantError.Error {
					t.Errorf("CalcHandler returned %v, but want %v", errorResponse.Error, tc.wantError.Error)
				}
			}
		})
	}
}

func TestCalcHandlerInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"expression": "2+2"`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	db, err := repository.InitDB("../db/calc.db")
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}

	handler := http.HandlerFunc(CalcHandler(repository.NewExpressionModel(db)))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned %v, but want %v", status, http.StatusBadRequest)
	}

	var errorResponse models.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to decode error response body: %v", err)
	}
	wantError := "Bad request"
	if errorResponse.Error != wantError {
		t.Errorf("Handler returned %v, but want %v", errorResponse.Error, wantError)
	}
}
