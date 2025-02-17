package models

import "errors"

var (
	ErrorDivisionByZero     = errors.New("division by zero is not allowed")
	ErrorEmptyBrackets      = errors.New("the brackets are empty")
	ErrorEmptyExpression    = errors.New("expression is empty")
	ErrorInvalidCharacter   = errors.New("invalid characters in expression")
	ErrorInvalidInput       = errors.New("expression is not valid")
	ErrorInvalidOperand     = errors.New("an invalid operand")
	ErrorInvalidRequestBody = errors.New("invalid request body")
	ErrorMissingOperand     = errors.New("missing operand")
	ErrorUnclosedBracket    = errors.New("the brackets in the expression are not consistent")
)

type ErrorResponse struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"error_message"`
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}

type Token struct {
	Value    string
	IsNumber bool
}

func NewToken(value string, isNumber bool) *Token {
	newToken := Token{
		Value:    value,
		IsNumber: isNumber,
	}
	return &newToken
}
