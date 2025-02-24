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

// GetTask забирает из базы таску для агента
func (e *ExpressionModel) GetTask() (*models.Task, error) {
	query := `
        SELECT id, expressionID, arg1, arg2, prev_task_id1, prev_task_id2, operation, status, result
        FROM tasks
        WHERE status = ?
        AND (prev_task_id1 = 0 OR prev_task_id1 IN (SELECT id FROM tasks WHERE status = ?))
        AND (prev_task_id2 = 0 OR prev_task_id2 IN (SELECT id FROM tasks WHERE status = ?))
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

	if err := e.LockTask(task.ID); err != nil {
		return nil, err
	}

	return &task, nil
}

// GetTaskStatus возвращает статус и ответ таски
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

// LockTask блокирует таску
func (e *ExpressionModel) LockTask(taskID int) error {
	_, err := e.DB.Exec("UPDATE tasks SET locked = 1 WHERE id = ?", taskID)
	return err
}

// UpdateTaskResult обновляет ответ таски в базе
func (e *ExpressionModel) UpdateTaskResult(taskID int, result float64) error {
	_, err := e.DB.Exec(
		"UPDATE tasks SET status = ?, result = ? WHERE id = ?",
		models.StatusResolved,
		result,
		taskID,
	)
	return err
}
