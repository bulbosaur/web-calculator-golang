package repository

import (
	"database/sql"
	"fmt"
	"log"
)

// InitDB выносит логику открытия БД в отдельную функцию
func InitDB(path string) (*sql.DB, error) {
	log.Printf("Database path: %s", path)

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("error when opening database: %v", err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS expressions (id INTEGER PRIMARY KEY AUTOINCREMENT, expression TEXT NOT NULL, status TEXT NOT NULL, result TEXT);`,
	)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error when connecting with database: %v", err)
	}

	log.Print("Successful connection to the database")

	return db, nil
}
