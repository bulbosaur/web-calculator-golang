package repository

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
)

// ExpressionModel обертывает пул подключения sql.DB
type ExpressionModel struct {
	DB *sql.DB
	mu sync.Mutex
}

// NewExpressionModel создает экземпляр ExpressionModel
func NewExpressionModel(db *sql.DB) *ExpressionModel {
	return &ExpressionModel{DB: db}
}

// AreAllTasksCompleted проверяет, все ли таски данного выражения выполнены
func (e *ExpressionModel) AreAllTasksCompleted(exprID int) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM tasks 
        WHERE expressionID = ? AND status != ?
    `
	var count int
	err := e.DB.QueryRow(query, exprID, models.StatusResolved).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check tasks completion: %v", err)
	}
	return count == 0, nil
}

// CalculateExpressionResult выбирает результаты всех тасок задачи и возвращает итоговый
func (e *ExpressionModel) CalculateExpressionResult(exprID int) (float64, error) {
	query := `
        SELECT result 
        FROM tasks 
        WHERE expressionID = ? AND status = ?
    `
	rows, err := e.DB.Query(query, exprID, models.StatusResolved)
	if err != nil {
		return 0, fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	var results []float64
	for rows.Next() {
		var result float64
		if err := rows.Scan(&result); err != nil {
			return 0, fmt.Errorf("failed to scan task result: %v", err)
		}
		results = append(results, result)
	}

	if len(results) == 0 {
		return 0, fmt.Errorf("no completed tasks found for expression ID %d", exprID)
	}

	return results[len(results)-1], nil
}

// Insert записывает мат выражение в таблицу БД
func (e *ExpressionModel) Insert(expression string) (int, error) {
	query := "INSERT INTO expressions (expression, status, result) VALUES (?, ?, ?)"

	result, err := e.DB.Exec(query, expression, models.StatusWait, "")
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorCreatingDatabaseRecord, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorReceivingID, err)
	}

	return int(id), nil
}

// GetExpression возвращает из базы данных соответствующее выражение
func (e *ExpressionModel) GetExpression(exprID int) (*models.Expression, error) {
	query := `
	SELECT id, status, result
	FROM expressions
	WHERE id = ?
	`
	var expr models.Expression

	err := e.DB.QueryRow(query, exprID).Scan(
		&expr.ID,
		&expr.Status,
		&expr.Result,
	)
	if err != nil {
		return nil, fmt.Errorf("fail to get expression ID-%d: %v", exprID, err)
	}

	return &expr, nil
}

// UpdateExpressionResult обновляет результат и статус выражения
func (e *ExpressionModel) UpdateExpressionResult(exprID int, result float64) error {
	query := `
        UPDATE expressions 
        SET result = ?, status = ? 
        WHERE id = ?
    `
	_, err := e.DB.Exec(query, result, models.StatusResolved, exprID)
	if err != nil {
		return fmt.Errorf("failed to update expression result: %v", err)
	}
	return nil
}

// UpdateStatus устанавливает актуальный статус выражения в БД
func (e *ExpressionModel) UpdateStatus(id int, status string) {
	query := "UPDATE expressions SET status = ? WHERE id = ?"

	_, err := e.DB.Exec(query, status, id)
	if err != nil {
		log.Println(err)
	}
}

// SetResult вносит в базу данных ответ на выражение
func (e *ExpressionModel) SetResult(id int, result float64) {
	query := "UPDATE expressions SET result = ? WHERE id = ?"

	_, err := e.DB.Exec(query, result, id)
	if err != nil {
		log.Println(err)
	}
}
