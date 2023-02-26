package repository

import (
	"database/sql"
	"fmt"
	"scrapper/internal/domain/entity"
)

const nodesTable = "nodes"

type Storage struct {
	db *sql.DB
}

func NewNodesRepository(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) InsertNode(node entity.Node) error {
	query := fmt.Sprintf("INSERT INTO %s(oid, name, sub_children, sub_nodes_total, description) VALUES(?,?,?,?,?)", nodesTable)
	_, err := s.db.Exec(query, node.OID, node.Name, node.SubChildren, node.SubNodesTotal, node.Description)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) InsertNodes(nodes []entity.Node) error {
	query := fmt.Sprintf("INSERT INTO %s(oid, name, sub_children, sub_nodes_total, description) VALUES(?,?,?,?,?)", nodesTable)

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		tx.Exec(query, node.OID, node.Name, node.SubChildren, node.SubNodesTotal, node.Description)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
