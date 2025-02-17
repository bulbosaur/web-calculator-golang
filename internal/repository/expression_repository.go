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
	query := "INSERT INTO expression (expression, status, result) VALUES (?, ?, ?)"
	result, err := e.DB.Exec(query, expression, "awaiting processing", "")

	if err != nil {
		log.Printf("something went wrong while creating a record in the database: %d", err)
		return 0, fmt.Errorf("%w: %v", models.ErrorCreatingDatabaseRecord, err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrorReceivingID, err)
	}

	return int(id), nil
}
