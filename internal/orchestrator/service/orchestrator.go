package orchestrator

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/spf13/viper"
)

// Calc вызывает токенизацию выражения, записывает его в RPN. а затем в параллельных горутинах подсчитывает значения выражений в скобках
func Calc(stringExpression string) (float64, error) {
	agentCalculations, err := processTasks(stringExpression)
	if err != nil {
		return 0, err
	}

	expression, err := tokenize(agentCalculations)
	if err != nil {
		return 0, err
	}

	if len(expression) == 0 {
		return 0, models.ErrorEmptyExpression
	}

	reversePolishNotation, err := toReversePolishNotation(expression)
	if err != nil {
		return 0, err
	}

	return evaluateRPN(reversePolishNotation)
}

func processTasks(expr string) (string, error) {
	re := regexp.MustCompile(`\(([^()]+)\)`)
	matches := re.FindAllStringSubmatch(expr, -1)
	if matches == nil {
		return expr, nil
	}

	type result struct {
		original string
		computed float64
		err      error
	}

	var wg sync.WaitGroup
	resultsChan := make(chan result, len(matches))

	for _, m := range matches {
		original := m[0]
		task := m[1]
		wg.Add(1)
		go func(task, original string) {
			defer wg.Done()
			computed, err := sendTaskToAgent(task)
			resultsChan <- result{original: original, computed: computed, err: err}
		}(task, original)
	}

	wg.Wait()
	close(resultsChan)

	newExpr := expr
	for res := range resultsChan {
		if res.err != nil {
			return "", res.err
		}
	}
	if strings.Contains(newExpr, "(") && strings.Contains(newExpr, ")") {
		return processTasks(newExpr)
	}
	return newExpr, nil
}

func sendTaskToAgent(task string) (float64, error) {
	host := viper.GetString("server.ORC_HOST")
	port := viper.GetString("server.ORC_PORT")
	taskURL := fmt.Sprintf("http://%s:%s/internal/task", host, port)
	resp, err := http.PostForm(taskURL,
		map[string][]string{
			"expression": {task},
		})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("agent returned status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	resultStr := strings.TrimSpace(string(bodyBytes))
	value, err := strconv.ParseFloat(resultStr, 64)
	if err != nil {
		return 0, errors.New("failed to parse agent response")
	}
	return value, nil
}
