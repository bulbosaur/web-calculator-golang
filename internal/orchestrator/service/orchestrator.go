package orchestrator

import (
	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/bulbosaur/web-calculator-golang/internal/repository"
)

// Calc вызывает токенизацию выражения, записывает его в RPN. а затем в параллельных горутинах подсчитывает значения выражений в скобках
func Calc(stringExpression string, id int, taskRepo *repository.ExpressionModel) error {
	expression, err := tokenize(stringExpression)
	if err != nil {
		return err
	}

	if len(expression) == 0 {
		return models.ErrorEmptyExpression
	}

	reversePolishNotation, err := toReversePolishNotation(expression)
	if err != nil {
		return err
	}

	parseRPN(reversePolishNotation, id, taskRepo)

	return nil
}

// func sendTaskToAgent(task *models.Task) (models.Token, error) {
// 	host := viper.GetString("server.ORC_HOST")
// 	port := viper.GetString("server.ORC_PORT")
// 	taskURL := fmt.Sprintf("http://%s:%s/internal/task", host, port)

// 	taskJSON, err := json.Marshal(task)
// 	if err != nil {
// 		return models.Token{}, fmt.Errorf("failed to marshal task: %v", err)
// 	}

// 	resp, err := http.Post(taskURL, "application/json", bytes.NewBuffer(taskJSON))
// 	if err != nil {
// 		return models.Token{}, fmt.Errorf("failed to send task to agent: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return models.Token{}, fmt.Errorf("agent returned status %s", resp.Status)
// 	}

// 	var resultToken models.Token
// 	if err := json.NewDecoder(resp.Body).Decode(&resultToken); err != nil {
// 		return models.Token{}, fmt.Errorf("failed to decode agent response: %v", err)
// 	}

// 	return resultToken, nil
// }

func NewTask(id int, arg1, arg2 float64, operation string) *models.Task {
	newTask := models.Task{
		ExpressionID: id,
		Arg1:         arg1,
		Arg2:         arg2,
		Operation:    operation,
	}
	return &newTask
}
