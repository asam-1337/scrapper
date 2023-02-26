package controller

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"math/rand"
	"net/http"
	"scrapper/internal/domain"
	"scrapper/internal/localErrors"
	"time"
)

const (
	UserAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 YaBrowser/23.1.2.987 Yowser/2.5 Safari/537.36"
	Accept         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	AcceptLanguage = "ru,en;q=0.9"
)

type Client struct {
	client   *http.Client
	svc      domain.NodesService
	attempts int
	bar      *progressbar.ProgressBar
}

func NewClient(svc domain.NodesService) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
		svc: svc,
		bar: progressbar.DefaultBytes(
			-1,
			"total nodes parsing",
		),
	}
}

func (c *Client) Do(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cant create req: %s", err.Error())
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", Accept)
	req.Header.Set("Accept-Language", AcceptLanguage)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cant do req (url: %s, err: %s)", url, err.Error())
	}

	return resp, nil
}

func (c *Client) Parse(url string) error {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	resp, err := c.Do(url)
	if err != nil {
		return err
	}

	c.bar.Add(1)

	nodes, err := c.svc.InsertNodes(resp)
	if err == localErrors.ErrNotFoundChildren {
		return nil
	}
	if err != nil {
		return fmt.Errorf("cant insert nodes: %s", err.Error())
	}

	for _, node := range nodes {
		url = "https://oidref.com" + node.OID
		retry(c.Parse, url)
	}

	return nil
}

func retry(f func(string) error, url string) {
	err := f(url)

	for err != nil {
		for i := 0; i < 10; i++ {
			err = f(url)
			if err != nil {
				log.Printf("url:%s, attemption: %d, error: %s", url, i, err.Error())
			}
		}
		time.Sleep(60 * time.Second)
	}
}

func (c *Client) StartParsing() {
	retry(c.Parse, "https://oidref.com/")
}
