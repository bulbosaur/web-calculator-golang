package orchestrator

import (
	"strconv"
	"sync"
	"time"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/spf13/viper"
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

func evaluateRPN(rpn []models.Token) (float64, error) {
	stackResultPolish := make([]float64, 0)
	var wg sync.WaitGroup
	results := make(chan float64, len(rpn))
	errors := make(chan error, len(rpn))

	for _, token := range rpn {
		if token.Value == "(" {
			wg.Add(1)
			go func(t models.Token) {
				defer wg.Done()
				result, err := evaluateSubExpression(t)
				if err != nil {
					errors <- err
					return
				}
				results <- result
			}(token)
		} else {
			floatNumber, err := strconv.ParseFloat(token.Value, 64)
			if err == nil {
				stackResultPolish = append(stackResultPolish, floatNumber)
			} else {
				if len(stackResultPolish) < 2 {
					return 0, models.ErrorInvalidInput
				}

				num1 := stackResultPolish[len(stackResultPolish)-1]
				num2 := stackResultPolish[len(stackResultPolish)-2]
				stackResultPolish = stackResultPolish[:len(stackResultPolish)-2]

				switch token.Value {
				case "+":
					time.Sleep(time.Duration(viper.GetInt("duration.TIME_ADDITION_MS")) * time.Millisecond)
					stackResultPolish = append(stackResultPolish, num1+num2)
				case "-":
					time.Sleep(time.Duration(viper.GetInt("duration.TIME_SUBTRACTION_MS")) * time.Millisecond)
					stackResultPolish = append(stackResultPolish, num2-num1)
				case "*":
					time.Sleep(time.Duration(viper.GetInt("duration.TIME_MULTIPLICATIONS_MS")) * time.Millisecond)
					stackResultPolish = append(stackResultPolish, num2*num1)
				case "/":
					time.Sleep(time.Duration(viper.GetInt("duration.TIME_DIVISIONS_MS")) * time.Millisecond)
					if num1 == 0 {
						return 0, models.ErrorDivisionByZero
					}
					stackResultPolish = append(stackResultPolish, num2/num1)
				}
			}
		}
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	for result := range results {
		stackResultPolish = append(stackResultPolish, result)
	}

	if len(errors) > 0 {
		return 0, <-errors
	}

	return stackResultPolish[0], nil
}

func evaluateSubExpression(token models.Token) (float64, error) {
	return Calc(token.Value)
}

func lastToken(tokens []models.Token) models.Token {
	return tokens[len(tokens)-1]
}
