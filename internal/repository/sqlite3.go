package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./migrations/nodes.db")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(15)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
