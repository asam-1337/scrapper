package domain

import (
	"net/http"
	"scrapper/internal/domain/entity"
)

type ParserService interface {
	ParseNodes(resp *http.Response) ([]entity.Node, error)
}
