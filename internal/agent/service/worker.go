package agent

import (
	"log"
	"time"
)

func worker(id int, orchestratorURL string) {
	interval := 5 * time.Second
	for {
		task, err := getTask(orchestratorURL)
		if err != nil {
			log.Printf("worker %d: task receiving error: %v", id, err)
			time.Sleep(interval)
			continue
		}

		if task == nil {
			log.Printf("Worker %d: no tasks available", id)
			time.Sleep(interval)
			continue
		}

		log.Printf("Worker %d: received task ID-%d", id, task.ID)
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
