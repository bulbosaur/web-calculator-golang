package models

// ErrorResponse - структура ответа, возвращаемого при ошибке вычислений
type ErrorResponse struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"error_message"`
}

// Request - структура запроса
type Request struct {
	Expression string `json:"expression"`
}

// Response - струтура ответа после успешного завершения программы
type Response struct {
	Result float64 `json:"result"`
}

// Token - структура токена, на которые разбивается исходное выражение
type Token struct {
	Value    string
	IsNumber bool
}
