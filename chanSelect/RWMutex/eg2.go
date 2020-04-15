package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var count int
var rw sync.RWMutex

func main() {
	ch := make(chan struct{}, 10)
	for i := 0; i < 5; i++ {
		go read(i, ch)
	}
	for i := 0; i < 5; i++ {
		go write(i, ch)
	}
	for i := 0; i < 10; i++ {
		<-ch
	}

}
func read(n int, ch chan struct{}) {
	rw.RLock()
	fmt.Printf("goroutine%d进入读操作。。。\n", n)
	v := count
	fmt.Printf("goroutine%d读取结束，值为：%d\n", n, v)
	rw.RUnlock()
	ch <- struct{}{}
}

func write(n int, ch chan struct{}) {
	rw.Lock()
	fmt.Printf("gouroutiine%d进入写操作。。。\n", n)
	count = rand.Intn(1000)
	fmt.Printf("goroutine%d写入结束，新值%d\n", n, count)
	rw.Unlock()
	ch <- struct{}{}
}
