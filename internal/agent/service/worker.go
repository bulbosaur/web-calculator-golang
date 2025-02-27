package agent

import (
	"log"
	"sync"
	"time"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
)

var Mu sync.Mutex

func worker(id int, orchestratorURL string) {
	interval := 1 * time.Second
	for {
		Mu.Lock()

		task, err := getTask(orchestratorURL)

		Mu.Unlock()

		if err != nil {
			log.Printf("worker %d: task receiving error: %v", id, err)
			time.Sleep(interval)
			continue
		}

		result, err := executeTask(orchestratorURL, task)
		if err == models.ErrorDivisionByZero {
			err = sendResult(orchestratorURL, task.ID, result, models.ErrorDivisionByZero)
		}

		if err != nil && task.ID != 0 {
			log.Printf("Worker %d: execution error task ID-%d: %v", id, task.ID, err)
			time.Sleep(interval)
			continue
		}

		if task.ID != 0 {
			err = sendResult(orchestratorURL, task.ID, result, nil)
			if err != nil {
				log.Printf("Worker %d: sending error task ID-%d: %v", id, task.ID, err)
			} else {
				log.Printf("Worker %d: success task ID-%d\nresult: %f", id, task.ID, result)
			}
		}

		time.Sleep(interval)
	}
}
