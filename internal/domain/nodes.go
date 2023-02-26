package domain

import (
	"net/http"
	"scrapper/internal/domain/entity"
)

type NodesRepository interface {
	InsertNode(node entity.Node) error
	InsertNodes(nodes []entity.Node) error
}

type NodesService interface {
	InsertNodes(resp *http.Response) ([]entity.Node, error)
}
