package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
)

// InsertTask записывает мат выражение в таблицу БД
func (e *ExpressionModel) InsertTask(task *models.Task) (int, error) {
	query := `
        INSERT INTO tasks (expressionID, arg1, arg2, prev_task_id1, prev_task_id2, operation, status, result)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	result, err := e.DB.Exec(
		query,
		task.ExpressionID,
		task.Arg1,
		task.Arg2,
		task.PrevTaskID1,
		task.PrevTaskID2,
		task.Operation,
		task.Status,
		task.Result,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get task ID: %v", err)
	}

	return int(id), nil
}

func (e *ExpressionModel) GetTask() (*models.Task, error) {
	query := `
        SELECT id, expressionID, arg1, arg2, prev_task_id1, prev_task_id2, operation, status, result
        FROM tasks
        WHERE status = ? AND (prev_task_id1 = 0 OR prev_task_id1 IS NULL OR prev_task_id1 IN (SELECT id FROM tasks WHERE status = ?))
          AND (prev_task_id2 = 0 OR prev_task_id2 IS NULL OR prev_task_id2 IN (SELECT id FROM tasks WHERE status = ?))
        LIMIT 1
    `

	var task models.Task
	err := e.DB.QueryRow(query, models.StatusWait, models.StatusResolved, models.StatusResolved).Scan(
		&task.ID,
		&task.ExpressionID,
		&task.Arg1,
		&task.Arg2,
		&task.PrevTaskID1,
		&task.PrevTaskID2,
		&task.Operation,
		&task.Status,
		&task.Result,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return &task, nil
}

func (e *ExpressionModel) GetTaskStatus(taskID int) (string, float64, error) {
	var status string
	var result float64

	err := e.DB.QueryRow("SELECT status, result FROM tasks WHERE id = ?", taskID).Scan(&status, &result)
	if err != nil {
		return "", 0, err
	}
	log.Printf("status %v", status)
	return status, result, nil
}
