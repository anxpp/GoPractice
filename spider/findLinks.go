package spider

import (
	"net/http"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

func Run(url string) {
	links, err := findLinks(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findLinks error: %v\n", err)
	}
	for _, link := range links {
		fmt.Println(link)
	}
}

// 下载网页
func findLinks(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("gettiong %s: %s", url, res.Status)
	}
	doc, err := html.Parse(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visit(&links, doc)
	return links, nil
}

// 查找超链接
func visit(links *[]string, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && !strings.Contains(a.Val, "javascript") {
				*links = append(*links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visit(links, c)
	}
}
