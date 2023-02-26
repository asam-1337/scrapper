package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"scrapper/internal/domain/entity"
	"scrapper/internal/localErrors"
	"strconv"
)

const (
	H3Exp    = `<h3>Children \(.*?\)</h3>`
	RowsExp  = `<tr>((?:.|\n)*?)</tr>`
	CellsExp = `(<td>(?:.|\n)*?"(?P<url>/.)*".*?</td>)(<td>(?P<name>.*?)</td>)(<td>(?P<subChildren>.*?)</td>)(<td>(?P<subNodesTotal>.*?)</td>)(<td>(?P<description>(?:.|\n)*?)</td>)`
)

type Parser struct {
	tableExp *regexp.Regexp
	h3Exp    *regexp.Regexp
	rowsExp  *regexp.Regexp
}

func NewParserService() *Parser {
	reTable, err := regexp.Compile(`<table>(.|\n)*?</table>`)
	if err != nil {
		log.Println(err)
		return nil
	}

	reH3, err := regexp.Compile(`<h3>Children.*?</h3>`)
	if err != nil {
		log.Println(err)
		return nil
	}

	reRows, err := regexp.Compile(`<tr>(.|\n)*?(<td>(.*?)"(?P<url>.*?)"(.*?)</td>)(.|\n)*?(<td>(?P<name>.*?)</td>)(.|\n)*?(<td>(?P<subChildren>.*?)</td>)(.|\n)*?(<td>(?P<subNodesTotal>.*?)</td>)(.|\n)*?(<td>(?P<description>.*?)</td>)(.|\n)*?</tr>`)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &Parser{
		tableExp: reTable,
		h3Exp:    reH3,
		rowsExp:  reRows,
	}
}

func (p *Parser) ParseNodes(resp *http.Response) ([]entity.Node, error) {
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cant read resp.Body: %s", err.Error())
	}

	url := resp.Request.URL.Path
	childrenHeader := p.h3Exp.FindString(string(html))
	if childrenHeader == "" && url != "/" {
		return nil, localErrors.ErrNotFoundChildren
	}

	nodes := make([]entity.Node, 0)
	table := p.tableExp.FindString(string(html))
	rows := p.rowsExp.FindAllStringSubmatch(table, -1)

	for _, row := range rows {
		n := entity.Node{}
		for i, groupName := range p.rowsExp.SubexpNames() {
			switch groupName {
			case "url":
				n.OID = row[i]
			case "name":
				n.Name = row[i]
			case "subChildren":
				sc, err := strconv.Atoi(row[i])
				if err != nil {
					return nil, err
				}
				n.SubChildren = sc
			case "subNodesTotal":
				snt, err := strconv.Atoi(row[i])
				if err != nil {
					return nil, err
				}
				n.SubNodesTotal = snt
			case "description":
				n.Description = row[i]
			}
		}
		nodes = append(nodes, n)
	}

	return nodes, nil
}
