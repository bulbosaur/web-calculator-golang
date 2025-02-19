package repository

import (
	"database/sql"
	"fmt"
	"log"
)

// InitDB выносит логику открытия БД в отдельную функцию
func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("error when opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error when connecting with database: %v", err)
	}

	log.Print("Successful connection to the database")

	return db, nil
}
