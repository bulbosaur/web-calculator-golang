package orchestrator

import (
	"unicode"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
)

func newToken(value string, isNumber bool) *models.Token {
	newToken := models.Token{
		Value:    value,
		IsNumber: isNumber,
	}
	return &newToken
}

func tokenize(expression string) ([]models.Token, error) {
	var (
		tokens []models.Token
		number string
		err    error
	)

	for i, symbol := range expression {
		if unicode.IsSpace(symbol) {
			if number != "" {
				tokens = append(tokens, *newToken(number, true))
			}
			continue
		}

		if unicode.IsDigit(symbol) {
			number += string(symbol)
			if i+1 == len(expression) || !unicode.IsDigit(rune(expression[i+1])) {
				tokens = append(tokens, *newToken(number, true))
				number = ""
			}
			continue
		}

		switch string(symbol) {
		case "+", "-", "/", "*", "(", ")":
			tokens = append(tokens, *newToken(string(symbol), false))
		default:
			err = models.ErrorInvalidCharacter
			return nil, err
		}

	}

	if number != "" {
		tokens = append(tokens, *newToken(number, true))
		number = ""
	}

	if !checkMissingOperand(tokens) {
		return nil, models.ErrorMissingOperand
	}

	if !checkEmptyBrackets(tokens) {
		return nil, models.ErrorEmptyBrackets
	}

	result := addMissingOperand(tokens)
	return result, nil
}

func checkMissingOperand(tokens []models.Token) bool {
	for i, token := range tokens {
		if i == len(tokens)-1 {
			break
		}
		if token.IsNumber && tokens[i+1].IsNumber {
			return false
		}
	}
	return true
}

func checkEmptyBrackets(tokens []models.Token) bool {
	for i, token := range tokens {
		if i == len(tokens)-1 {
			break
		}
		if token.Value == "(" && tokens[i+1].Value == ")" {
			return false
		}
	}
	return true
}

func addMissingOperand(expression []models.Token) []models.Token {
	var result []models.Token

	for i, token := range expression {
		result = append(result, token)

		if i+1 < len(expression) {
			if (token.IsNumber || token.Value == ")") && expression[i+1].Value == "(" {
				result = append(result, models.Token{Value: "*", IsNumber: false})
			}
			if token.Value == ")" && expression[i+1].IsNumber {
				result = append(result, models.Token{Value: "*", IsNumber: false})
			}
		}
	}

	return result
}
