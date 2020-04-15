package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter = struct {
		sync.Map
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}
)

func main() {
	go func() {
		for {
			counter.RLock()
			n := counter.m["some_key"]
			counter.RUnlock()
			fmt.Println(n)
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for {
			counter.Lock()
			counter.m["some_key"]++
			counter.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(10 * time.Second)
}
