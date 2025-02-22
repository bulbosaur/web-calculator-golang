package models

var (
	// StatusFailed указывает, что выражение не решено. Причиной может быть его некорректность
	StatusFailed = "failed"

	// StatusResolved указывает в БД, что результат выражения подсчитан успешно
	StatusResolved = "successfully done"

	// StatusWait указывает на те выражения в БД, результат которых еще не подсчитан
	StatusWait = "awaiting processing"
)

// ErrorResponse - структура ответа, возвращаемого при ошибке вычислений
type ErrorResponse struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"error_message"`
}

// RegisteredExpression - структура ответа, возвращаемого при регистрации выражения в оркестраторе
type RegisteredExpression struct {
	Id int `json:"id"`
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
