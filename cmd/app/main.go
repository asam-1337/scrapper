package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const (
	UserAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 YaBrowser/23.1.2.987 Yowser/2.5 Safari/537.36"
	Accept         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	AcceptEncoding = "gzip, deflate, br"
	AcceptLanguage = "ru,en;q=0.9"
)

func main() {
	c := http.Client{}
	req, err := http.NewRequest("GET", "https://oidref.com/0.2", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", Accept)
	//req.Header.Set("Accept-Encoding", AcceptEncoding)
	req.Header.Set("Accept-Language", AcceptLanguage)

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Status)
	//fmt.Println(string(b))
	re1, err := regexp.Compile(`Children.*?Brothers`)
	re2, err := regexp.Compile(`<tr>.*?</tr>`)
	re3, err := regexp.Compile(`<td>.*?<td>`)
	table := re1.FindAllString(string(b), -1)
	fmt.Println(table)
	rows := re2.FindAllString(table[0], -1)
	for i, v := range rows {
		if i == 0 {
			continue
		}
		data := re3.FindAllString(v, -1)
		for _, d := range data {
			fmt.Print(d)
		}
		fmt.Print("\n")
	}

}
