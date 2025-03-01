package agent

import (
	"log"
	"sync"
	"time"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
)

// Mu - мьютекс в рамках микросервиса данного агента
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

		result, errorMessage, err := executeTask(orchestratorURL, task)
		if err != nil && task.ID != 0 {

			log.Printf("Worker %d: execution error task ID-%d: %v", id, task.ID, err)
			time.Sleep(interval)
			continue
		}

		if task.ID != 0 {
			err = sendResult(orchestratorURL, task.ID, result, errorMessage)
			if err != nil {
				log.Printf("Worker %d: sending error task ID-%d: %v", id, task.ID, err)
			} else {
				log.Printf("Worker %d: success task ID-%d\nresult: %f", id, task.ID, result)
			}
		}

		time.Sleep(interval)
	}
}
