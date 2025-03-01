package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// InitDB открывает соединение с базой и создаёт необходимые таблицы
func InitDB(path string) (*sql.DB, error) {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		log.Printf("Created directory: %s", dir)
	}

	log.Printf("Database path: %s", path)

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("error when opening database: %v", err)
	}

	createExpressions := `
 CREATE TABLE IF NOT EXISTS expressions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  expression TEXT NOT NULL,
  status TEXT NOT NULL,
  result TEXT,
  error_message TEXT
 );`
	_, err = db.Exec(createExpressions)
	if err != nil {
		return nil, fmt.Errorf("error creating expressions table: %v", err)
	}

	createTasks := `
 CREATE TABLE IF NOT EXISTS tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  expressionID INTEGER NOT NULL,
  arg1 TEXT NOT NULL,
  arg2 TEXT NOT NULL,
  prev_task_id1 INTEGER DEFAULT 0,
  prev_task_id2 INTEGER DEFAULT 0,
  operation TEXT NOT NULL,
  status TEXT,
  result FLOAT,
  error_message TEXT,
  FOREIGN KEY(expressionID) REFERENCES expressions(id)
 );`
	_, err = db.Exec(createTasks)
	if err != nil {
		return nil, fmt.Errorf("error creating tasks table: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error when connecting with database: %v", err)
	}

	log.Print("Successful connection to the database")
	return db, nil
}
