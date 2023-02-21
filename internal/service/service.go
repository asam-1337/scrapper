package service

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"scrapper/internal/repository"
)

const (
	TableExp = ``
	TagsExp  = `<tr><td><a href=".*?">`
	LinksExp = `/[0-9.]+?`
)

type Interface interface {
	ParseLinks(resp *http.Response) ([]string, error)
}

type ParserService struct {
	repo     repository.Repository
	tagsExp  *regexp.Regexp
	linksExp *regexp.Regexp
}

func NewService(repo *repository.Repository) *ParserService {
	reTags, err := regexp.Compile(TagsExp)
	if err != nil {
		log.Println(err)
		return nil
	}

	reLinks, err := regexp.Compile(LinksExp)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &ParserService{
		repo:     repo,
		tagsExp:  reTags,
		linksExp: reLinks,
	}
}

func (s *ParserService) ParseLinks(resp *http.Response) ([]string, error) {
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	links := s.tagsExp.FindAllString(string(b), -1)
	for i, val := range links {
		links[i] = s.linksExp.FindString(val)
	}

	return links, nil
}

//func (s *ParserService) ParseDescription(resp *http.Response) (string, error) {
//	defer resp.Body.Close()
//	b, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Println(err)
//		return "", err
//	}
//
//	desc := s.
//}
