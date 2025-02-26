package agent

import (
	"log"
	"sync"
	"time"
)

var Mu sync.Mutex

func worker(id int, orchestratorURL string) {
	interval := 5 * time.Second
	for {
		Mu.Lock()

		task, err := getTask(orchestratorURL)
		if err != nil {
			log.Printf("worker %d: task receiving error: %v", id, err)
			time.Sleep(interval)
			continue
		}

		Mu.Unlock()

		result, err := executeTask(orchestratorURL, task)
		if err != nil {
			log.Printf("Worker %d: execution error task ID-%d: %v", id, task.ID, err)
			time.Sleep(interval)
			continue
		}

		err = sendResult(orchestratorURL, task.ID, result)
		if err != nil {
			log.Printf("Worker %d: sending error task ID-%d: %v", id, task.ID, err)
		} else {
			log.Printf("Worker %d: success task ID-%d\nresult: %f", id, task.ID, result)
		}

		time.Sleep(interval)
	}
}
