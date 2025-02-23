package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
	"github.com/spf13/viper"
)

// RunAgent запускает агента
func RunAgent() {
	orc_host := viper.GetString("server.ORC_HOST")
	orc_port := viper.GetString("server.ORC_PORT")
	orchestratorURL := fmt.Sprintf("http://%s:%s", orc_host, orc_port)

	workers := viper.GetInt("COMPUTING_POWER")
	if workers <= 0 {
		workers = 1
	}

	for i := 0; i < workers; i++ {
		go worker(i, orchestratorURL)
	}

	select {}
}

func worker(id int, orchestratorURL string) {
	interval := 10 * time.Second
	for {
		task, err := getTask(orchestratorURL)
		if err != nil {
			log.Printf("worker %d: task receiving error: %v", id, err)
			time.Sleep(interval)
			continue
		}

		log.Printf("Worker %d: receive task ID-%d", id, task.Id)
		result, err := executeTask(task)
		if err != nil {
			log.Printf("Worker %d: execution error task ID-%d: %v", id, task.Id, err)
			time.Sleep(interval)
			continue
		}

		err = sendResult(orchestratorURL, task.Id, result)
		if err != nil {
			log.Printf("Worker %d: sending error task ID-%d: %v", id, task.Id, err)
		} else {
			log.Printf("Worker %d: success task ID-%d\nresult: %f", id, task.Id, result)
		}

		time.Sleep(interval)
	}
}

func getTask(orchestratorURL string) (*models.Task, error) {
	resp, err := http.Get(orchestratorURL + "/internal/task")
	log.Println(orchestratorURL + "/internal/task")
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("orchestrator returned status code %d", resp.StatusCode)
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode task: %w", err)
	}

	return &task, nil
}

func executeTask(task *models.Task) (float64, error) {
	switch task.Operation {
	case "+":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_ADDITION_MS")) * time.Millisecond)
		return task.Arg1 + task.Arg2, nil
	case "-":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_SUBTRACTION_MS")) * time.Millisecond)
		return task.Arg1 - task.Arg2, nil
	case "*":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_MULTIPLICATIONS_MS")) * time.Millisecond)
		return task.Arg1 * task.Arg2, nil
	case "/":
		time.Sleep(time.Duration(viper.GetInt("duration.TIME_DIVISIONS_MS")) * time.Millisecond)
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return task.Arg1 / task.Arg2, nil
	default:
		return 0, models.ErrorInvalidOperand
	}
}

func sendResult(orchestratorURL string, taskID int, result float64) error {
	payload, err := json.Marshal(map[string]interface{}{
		"id":     taskID,
		"result": result,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	resp, err := http.Post(orchestratorURL+"/internal/task", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("orchestrator returned status code %d", resp.StatusCode)
	}

	return nil
}
