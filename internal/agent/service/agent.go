package agent

import (
	"time"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/spf13/viper"
)

// ExecuteOperation выполняет одну операцию
func ExecuteOperation(arg1, arg2 float64, operation string) (float64, error) {
	switch operation {
	case "+":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_ADDITION_MS")) * time.Millisecond)
		return arg1 + arg2, nil
	case "-":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_SUBTRACTION_MS")) * time.Millisecond)
		return arg1 - arg2, nil
	case "*":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_MULTIPLICATIONS_MS")) * time.Millisecond)
		return arg1 * arg2, nil
	case "/":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_DIVISIONS_MS")) * time.Millisecond)
		if arg2 == 0 {
			return 0, models.ErrorDivisionByZero
		}
		return arg1 / arg2, nil
	default:
		return 0, models.ErrorInvalidInput
	}
}
