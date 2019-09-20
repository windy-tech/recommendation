package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func PrepareDatabase() (*sql.DB, error) {
	database, _ := sql.Open("sqlite3", "./windy.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, name TEXT, catagory TEXT, confidence FLOAT, content TEXT)")
	_, err := statement.Exec()
	if err != nil {
		return nil, err
	}
	return database, nil
}
