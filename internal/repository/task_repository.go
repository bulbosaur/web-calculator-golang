package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
)

// InsertTask записывает мат выражение в таблицу БД
func (e *ExpressionModel) InsertTask(task *models.Task, exprId int) (int, error) {
	query := "INSERT INTO tasks (expressionID, arg1, arg2, operation, operation_time, status, result) VALUES (?, ?, ?, ?, ?, ?, ?)"

	result, err := e.DB.Exec(query, exprId, task.Arg1, task.Arg2, task.Operation, 0, models.StatusWait, 0.0)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorCreatingDatabaseRecord, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorReceivingID, err)
	}

	return int(id), nil
}

func (e *ExpressionModel) GetTask() (*models.Task, error) {
	query := `
		SELECT id, expressionID, arg1, arg2, operation, operation_time, status, result
		FROM tasks
		WHERE status = ?
		LIMIT 1
	`

	var task models.Task

	err := e.DB.QueryRow(query, models.StatusWait).Scan(
		&task.Id,
		&task.ExpressionId,
		&task.Arg1,
		&task.Arg2,
		&task.Operation,
		&task.OperationTime,
		&task.Status,
		&task.Result,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%v", err)
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
