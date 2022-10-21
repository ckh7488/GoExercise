package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"goRegex/util"
)

var webSiteName string = "https://go.dev"

// var webSiteName string = "https://www.naver.com"
var maxdepth = 4

// func main() {
// 	resp, _ := http.Get(webSiteName)
// 	b, _ := io.ReadAll(resp.Body)
// 	// fmt.Println("%s", string(b))
// 	util.RetRegex(string(b))
// }

type SafeCounter struct {
	mu           sync.Mutex
	internalLink map[string]bool
	externalLink map[string]bool
}

func (c *SafeCounter) appendInternal(key string) {
	c.mu.Lock()
	c.internalLink[key] = true
	c.mu.Unlock()
}

func (c *SafeCounter) appendExternal(key string) {
	c.mu.Lock()
	c.externalLink[key] = true
	c.mu.Unlock()
}

func (c *SafeCounter) find(key string, isInternal bool) bool {
	if isInternal {
		if c.internalLink[key] {
			return true
		}
		return false
	}
	if c.externalLink[key] {
		return true
	}
	return false
}

func main() {
	num := int32(0)
	numGoroutine := &num
	done := make(chan bool)
	fi, _ := os.Create("test")
	sc := SafeCounter{internalLink: make(map[string]bool), externalLink: make(map[string]bool)}

	var exploreWeb func(*SafeCounter, string, int)
	exploreWeb = func(sc *SafeCounter, url string, depth int) {
		atomic.AddInt32(numGoroutine, 1)
		if depth == maxdepth {
			fmt.Println("depth is maxDepth")
			atomic.AddInt32(numGoroutine, -1)
			if *numGoroutine == 0 {
				done <- true
			}
			return
			// fmt.Println("%q", sc.internalLink)
			// fmt.Println("%q", sc.externalLink)
		}
		c := http.Client{Timeout: time.Duration(1) * time.Second * 2}
		resp, err := c.Get(webSiteName)
		if err != nil {
			fmt.Println("client Get error, url: %v err : %v", url, err)
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)
		strArr, _ := util.RetAllLinks(body)
		fmt.Println(len(strArr))
		// fmt.Println("%v", strArr)
		// fmt.Println("hgg")
		// fmt.Println("RetAllLinks ends")
		for _, link := range strArr {
			if len(link) > 4 && link[0:4] == "http" {
				sc.appendExternal(link)
				continue
			}
			if depth < maxdepth {
				nextDir := webSiteName + link
				// fmt.Println(sc.find(nextDir, true), nextDir)
				if !sc.find(nextDir, true) {
					sc.appendInternal(nextDir)
					go exploreWeb(sc, nextDir, depth+1)
				}
			}
		}
		atomic.AddInt32(numGoroutine, -1)
		if *numGoroutine == 0 {
			done <- true
		}
	}

	exploreWeb(&sc, webSiteName, 1)
	fmt.Println("exploreStart")
	<-done
	fmt.Println("exploreEnd")
	for i, _ := range sc.internalLink {
		str := fmt.Sprintf("link name : %v \n", i)
		fi.Write([]byte(str))
	}
	for i, _ := range sc.externalLink {
		str := fmt.Sprintf("external link name : %v \n", i)
		fi.Write([]byte(str))
	}
}
