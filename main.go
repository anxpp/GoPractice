package main

import (
	"./spider"
	"os"
	"fmt"
	"math/rand"
	"time"
	"net/http"
)

func main() {
	slice := spider.Rank("https://blog.csdn.net/anxpp", "https://blog.csdn.net/anxpp/article/details/ ")
	for _, value := range slice {
		fmt.Fprintf(os.Stderr, "%-3d %s \n", value.Rank, value.Url)
	}
	for {
		time.Sleep(time.Second * (time.Duration)(rand.Intn(3)+3))
		url := slice[rand.Intn(3)].Url
		fmt.Fprintf(os.Stdout, "gettiong %s \n", url)
		_, _ = http.Get(url)
	}
}
