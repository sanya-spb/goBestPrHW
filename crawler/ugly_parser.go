package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

// парсим страницу
func parse(url string) (*html.Node, error) {
	// что здесь должно быть вместо http.Get? :)
	// r, err := http.Get(url)
	client := http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if r, err := client.Get(url); err != nil {
		return nil, fmt.Errorf("can't get page")
	} else {
		if b, err := html.Parse(r.Body); err != nil {
			return nil, fmt.Errorf("can't parse page")
		} else {
			return b, err
		}
	}
}

// ищем заголовок на странице
func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

// ищем все ссылки на страницы. Используем мапку чтобы избежать дубликатов
func pageLinks(links map[string]struct{}, n *html.Node) map[string]struct{} {
	if links == nil {
		links = make(map[string]struct{})
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}

			// костылик для простоты
			if _, ok := links[a.Val]; !ok && len(a.Val) > 2 && a.Val[:2] == "//" {
				links["http://"+a.Val[2:]] = struct{}{}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = pageLinks(links, c)
	}
	return links
}
