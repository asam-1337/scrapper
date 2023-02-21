package controller

import (
	"log"
	"net/http"
	"scrapper/internal/localErrors"
	"scrapper/internal/repository"
	"scrapper/internal/service"
	"time"
)

const (
	UserAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 YaBrowser/23.1.2.987 Yowser/2.5 Safari/537.36"
	Accept         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	AcceptEncoding = "gzip, deflate, br"
	AcceptLanguage = "ru,en;q=0.9"
)

type Client struct {
	client *http.Client
	svc    service.Interface
}

func NewClient(repo *repository.Repository) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		svc: service.NewService(repo),
	}
}

func (c *Client) Do(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", Accept)
	//req.Header.Set("Accept-Encoding", AcceptEncoding)
	req.Header.Set("Accept-Language", AcceptLanguage)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) Parse(url string) error {
	resp, err := c.Do(url)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	links, err := c.svc.ParseLinks(resp)
	if len(links) == 0 {
		return localErrors.ErrNotFoundChildren
	}

	for _, link := range links {
		err = c.Parse(link)
		if err == localErrors.ErrNotFoundChildren {
			return nil
		}
	}

	return nil
}
