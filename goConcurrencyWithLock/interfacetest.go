package main

import (
	"errors"
	"fmt"
	"sync"
)

type lockedVar struct {
	mu          sync.Mutex
	criticalVar map[string]string
}

func (lv *lockedVar) addMap(key string, value string) {
	lv.mu.Lock()
	if v, ok := lv.criticalVar[key]; ok {
		fmt.Printf("value %v already saved. \n", v)
		lv.mu.Unlock()
		return
	}
	lv.criticalVar[key] = value
	lv.mu.Unlock()
}

func (lv *lockedVar) queryMap(key string) (string, error) {
	lv.mu.Lock()
	keyStr, ok := lv.criticalVar[key]
	lv.mu.Unlock()
	if ok {
		return keyStr, nil
	}
	return "", errors.New("no matching with this key")
}

func main() {
	tlv := lockedVar{criticalVar: make(map[string]string)}
	c := make(chan int)
	var addNum func(int, int, chan int)
	addNum = func(i int, mult int, c chan int) {
		for k := 0; k < i; k++ {
			idx := fmt.Sprint(k)
			val := fmt.Sprintf("%v", mult*k)
			tlv.addMap(idx, val)
		}
		c <- 1
	}

	go addNum(15, 1, c)
	go addNum(45, 100, c)
	go addNum(30, 10, c)
	<-c
	<-c
	<-c
	for idx, val := range tlv.criticalVar {
		fmt.Printf("%v : %v\n", idx, val)
	}
}
