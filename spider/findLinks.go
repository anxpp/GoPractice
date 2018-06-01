package spider

import (
	"net/http"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

func Run(url string) {
	links, images, err := findLinks(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findLinks error: %v\n", err)
	}
	fmt.Println("超链接如下：")
	printFun := func(links []string) {
		for _, link := range links {
			fmt.Println(link)
		}
	}
	printFun(links)
	fmt.Fprintf(os.Stderr, "图片数量：%d\n", len(images))
}

// 下载网页
func findLinks(url string) ([]string, []string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, nil, fmt.Errorf("gettiong %s: %s", url, res.Status)
	}
	doc, err := html.Parse(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	var images []string
	visit(&links, &images, doc)
	return links, images, nil
}

// 查找超链接
func visit(links *[]string, images *[]string, n *html.Node) {
	if n == nil {
		return
	}
	// 判断链接
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && !strings.Contains(a.Val, "javascript") {
				*links = append(*links, a.Val)
			}
		}
	}
	// 图片
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				*images = append(*images, a.Val)
			}
		}
	}
	// 广度优先
	visit(links, images, n.NextSibling)
	visit(links, images, n.FirstChild)
}
