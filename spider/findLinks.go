package spider

import (
	"net/http"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
	"sort"
	"time"
)

type Item struct {
	Url  string
	Rank int
}
type ItemList []Item

func (p ItemList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ItemList) Len() int           { return len(p) }
func (p ItemList) Less(i, j int) bool { return p[i].Rank < p[j].Rank }

var BaseUrl = ""

func Rank(url string, baseUrl string) ItemList {
	BaseUrl = baseUrl
	// key  ：url
	// value：[0]=是否抓取过，[1]=被其他网页引用过的次数
	all := make(map[string]*[2]int)
	Run(url, &all)
	var list ItemList
	for key, value := range all {
		list = append(list, Item{Url: key, Rank: value[1]})
	}
	sort.Sort(list)
	sort.Reverse(list)
	return list
}

func Run(url string, all *map[string]*[2]int) {
	// 已解析过就跳过
	info := (*all)[url]
	// 网页信息分配内存
	if info == nil {
		info = &([2]int{0, 0})
		(*all)[url] = info
	}
	// 如果已经下载过该网页则跳过
	if (*info)[0] == 1 {
		return
	}
	(*info)[0] = 1
	links, _, err := findLinks(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findLinks error: %v\n", err)
	}
	for _, link := range links {
		if strings.Contains(link, BaseUrl) {
			info := (*all)[link]
			if info == nil {
				info = &([2]int{0, 0})
			}
			(*info)[1]++
			Run(link, all)
		}
	}
	//fmt.Println("超链接如下：")
	//printFun := func(links []string) {
	//	for _, link := range links {
	//		fmt.Println(link)
	//	}
	//}
	//printFun(links)
	//fmt.Fprintf(os.Stderr, "图片数量：%d\n", len(images))
}

// 下载网页
func findLinks(url string) ([]string, []string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("gettiong %s: %s", url, res.Status)
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	// 查找内容
	var visit func(links *[]string, images *[]string, n *html.Node)
	visit = func(links *[]string, images *[]string, n *html.Node) {
		time.Sleep(time.Second)
		if n == nil {
			return
		}
		// 判断链接
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					url, err := res.Request.URL.Parse(a.Val)
					if err != nil {
						continue
					}
					*links = append(*links, url.String())
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
	var links []string
	var images []string
	visit(&links, &images, doc)
	return links, images, nil
}
