package orchestrator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/spf13/viper"
)

// Calc вызывает токенизацию выражения, записывает его в RPN. а затем в параллельных горутинах подсчитывает значения выражений в скобках
func Calc(stringExpression string) (float64, error) {
	expression, err := tokenize(stringExpression)
	if err != nil {
		return 0, err
	}

	if len(expression) == 0 {
		return 0, models.ErrorEmptyExpression
	}

	simpleExpression, err := taskSelection(expression)
	if err != nil {
		return 0, err
	}

	reversePolishNotation, err := toReversePolishNotation(simpleExpression)
	if err != nil {
		return 0, err
	}

	return evaluateRPN(reversePolishNotation)
}

func taskSelection(expression []models.Token) ([]models.Token, error) {
	var (
		simpleExpression []models.Token
		task             []models.Token
		inBrackets       bool
	)

	for _, token := range expression {
		if token.Value == "(" {
			inBrackets = true
			task = []models.Token{}
			continue
		}

		if token.Value == ")" {
			inBrackets = false

			for len(task) >= 3 {
				subTask := task[:3]
				resultToken, err := sendTaskToAgent(subTask)
				if err != nil {
					return nil, fmt.Errorf("failed to send subtask to agent: %v", err)
				}

				task = append([]models.Token{resultToken}, task[3:]...)
			}

			simpleExpression = append(simpleExpression, token)
			continue
		}

		if inBrackets {
			task = append(task, token)
		} else {
			simpleExpression = append(simpleExpression, token)
		}
	}

	return simpleExpression, nil
}

func sendTaskToAgent(task []models.Token) (models.Token, error) {
	host := viper.GetString("server.ORC_HOST")
	port := viper.GetString("server.ORC_PORT")
	taskURL := fmt.Sprintf("http://%s:%s/internal/task", host, port)

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return models.Token{}, fmt.Errorf("failed to marshal task: %v", err)
	}

	resp, err := http.Post(taskURL, "application/json", bytes.NewBuffer(taskJSON))
	if err != nil {
		return models.Token{}, fmt.Errorf("failed to send task to agent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Token{}, fmt.Errorf("agent returned status %s", resp.Status)
	}

	var resultToken models.Token
	if err := json.NewDecoder(resp.Body).Decode(&resultToken); err != nil {
		return models.Token{}, fmt.Errorf("failed to decode agent response: %v", err)
	}

	return resultToken, nil
}

func NewTask(id int, arg1, arg2 float64, operation string, operationTime int) *models.Task {
	newTask := models.Task{
		Id:            id,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operation,
		OperationTime: operationTime,
	}
	return &newTask
}
