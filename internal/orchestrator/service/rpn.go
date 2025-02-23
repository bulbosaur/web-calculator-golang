package orchestrator

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
)

func toReversePolishNotation(expression []models.Token) ([]models.Token, error) {
	priority := map[string]int{
		"(": 0,
		")": 1,
		"+": 2,
		"-": 2,
		"*": 3,
		"/": 3,
	}
	stack := []models.Token{}
	reversePolishNotation := []models.Token{}

	for _, token := range expression {
		if _, ok := priority[token.Value]; ok {
			if token.Value == ")" {
				for i := len(stack) - 1; i >= 0 && stack[i].Value != "("; i-- {
					reversePolishNotation = append(reversePolishNotation, lastToken(stack))
					stack = stack[:len(stack)-1]
				}

				if len(stack) > 0 && lastToken(stack).Value == "(" {
					stack = stack[:len(stack)-1]
				} else {
					return nil, models.ErrorUnclosedBracket
				}
				continue
			}

			for len(stack) > 0 && priority[lastToken(stack).Value] >= priority[token.Value] && token.Value != "(" {
				reversePolishNotation = append(reversePolishNotation, lastToken(stack))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)

		} else if token.IsNumber {
			reversePolishNotation = append(reversePolishNotation, token)
		} else {
			return nil, models.ErrorInvalidInput
		}
	}

	for len(stack) > 0 {

		reversePolishNotation = append(reversePolishNotation, lastToken(stack))
		stack = stack[:len(stack)-1]
	}
	return reversePolishNotation, nil
}

func parseRPN(expression []models.Token, id int, taskRepo *repository.ExpressionModel) error {
	var (
		stack []float64
	)

	for _, token := range expression {
		if token.IsNumber {
			value, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				log.Println(err)
			}
			stack = append(stack, value)
		} else {
			if len(stack) < 2 {
				return fmt.Errorf("there are not enough operands for the operation %s", token)
			}

			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			task := NewTask(id, arg1, arg2, token.Value)

			taskId, err := taskRepo.InsertTask(task, id)
			if err != nil {
				log.Printf("something went wrong while creating a record in the database. %v", err)
			}

			log.Printf("Task ID- %d (expression ID-%d) has been registered", taskId, id)
			for {
				taskStatus, taskResult, err := taskRepo.GetTaskStatus(taskId)
				if err != nil {
					log.Printf("failed to get task status: %v", err)
					return err
				}

				if taskStatus == models.StatusResolved {
					stack = append(stack, taskResult)
					break
				}

				time.Sleep(1 * time.Second)
			}
		}
	}

	if len(stack) != 1 {
		return models.ErrorInvalidInput
	}

	return nil
}

func lastToken(tokens []models.Token) models.Token {
	return tokens[len(tokens)-1]
}
