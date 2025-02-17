package models

import "errors"

var (
	//ErrorDivisionByZero - ошибка деления на ноль
	ErrorDivisionByZero = errors.New("division by zero is not allowed")

	// ErrorEmptyBrackets - пустые скобочки
	ErrorEmptyBrackets = errors.New("the brackets are empty")

	// ErrorEmptyExpression - пустое выражение
	ErrorEmptyExpression = errors.New("expression is empty")

	//ErrorInvalidCharacter - запрещенные символы в выражении
	ErrorInvalidCharacter = errors.New("invalid characters in expression")

	// ErrorInvalidInput - невалидное выражение
	ErrorInvalidInput = errors.New("expression is not valid")

	// ErrorInvalidOperand - ошибка при введении операнда
	ErrorInvalidOperand = errors.New("an invalid operand")

	// ErrorInvalidRequestBody - ошибка тела запроса
	ErrorInvalidRequestBody = errors.New("invalid request body")

	// ErrorMissingOperand - пропущенный операнд
	ErrorMissingOperand = errors.New("missing operand")

	// ErrorUnclosedBracket - скобочки не согласованы
	ErrorUnclosedBracket = errors.New("the brackets in the expression are not consistent")
)
