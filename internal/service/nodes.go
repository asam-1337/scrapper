package service

import (
	"net/http"
	"scrapper/internal/domain"
	"scrapper/internal/domain/entity"
)

type Nodes struct {
	repo   domain.NodesRepository
	parser domain.ParserService
}

func NewNodesService(repo domain.NodesRepository, parser domain.ParserService) *Nodes {
	return &Nodes{
		repo:   repo,
		parser: parser,
	}
}

func (s *Nodes) InsertNodes(resp *http.Response) ([]entity.Node, error) {
	nodes, err := s.parser.ParseNodes(resp)
	if err != nil {
		return nodes, err
	}

	err = s.repo.InsertNodes(nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
