package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bulbosaur/web-calculator-golang/internal/models"
)

// ExpressionModel обертывает пул подключения sql.DB
type ExpressionModel struct {
	DB *sql.DB
}

// NewExpressionModel создает экземпляр ExpressionModel
func NewExpressionModel(db *sql.DB) *ExpressionModel {
	return &ExpressionModel{DB: db}
}

// Insert записывает мат выражение в таблицу БД
func (e *ExpressionModel) Insert(expression string) (int, error) {
	query := "INSERT INTO expressions (expression, status, result) VALUES (?, ?, ?)"

	result, err := e.DB.Exec(query, expression, models.StatusResolved, "")
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorCreatingDatabaseRecord, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorReceivingID, err)
	}

	return int(id), nil
}

func (e *ExpressionModel) UpdateStatus(id int, status string) {
	query := "UPDATE expressions SET status = ? WHERE id = ?"

	_, err := e.DB.Exec(query, models.StatusResolved, id)
	if err != nil {
		log.Fatal(err)
	}
}
