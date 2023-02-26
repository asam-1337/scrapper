package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func NewSqliteDb() (*sql.DB, error) {
	_, err := os.Create("./migrations/nodes.db")
	if err != nil {
		return nil, fmt.Errorf("cant create file: %s", err.Error())
	}

	db, err := sql.Open("sqlite3", "./migrations/nodes.db")
	if err != nil {
		return nil, fmt.Errorf("cant open db conn: %s", err.Error())
	}

	createQuery := `CREATE TABLE IF NOT EXISTS nodes (
		 oid TEXT PRIMARY KEY ,
		 name TEXT,
		 sub_children INTEGER,
		 sub_nodes_total INTEGER,
		 description TEXT
	);`

	_, err = db.Exec(createQuery)
	if err != nil {
		return nil, fmt.Errorf("cant create table nodes: %s", err.Error())
	}

	db.SetMaxOpenConns(15)
	return db, nil
}
