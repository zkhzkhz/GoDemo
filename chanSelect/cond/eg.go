package main

import (
	"fmt"
	"sync"
	"time"
)

var count int = 4

func main() {
	ch := make(chan struct{}, 5)

	var l sync.Mutex
	cond := sync.NewCond(&l)

	for i := 0; i < 5; i++ {
		go func(i int) {
			cond.L.Lock()
			defer func() {
				cond.L.Unlock()
				ch <- struct{}{}
			}()

			for count > i {
				cond.Wait()
				fmt.Printf("收到一个通知goroutine%d\n", i)
			}
			fmt.Printf("gouroutine%d执行结束\n", i)
		}(i)
	}

	time.Sleep(time.Millisecond * 20)
	fmt.Println("broadcast...")
	cond.L.Lock()
	count -= 1
	cond.Broadcast()
	cond.L.Unlock()

	time.Sleep(time.Second)
	fmt.Println("signal...")
	cond.L.Lock()
	count -= 2
	cond.Signal()
	cond.L.Unlock()

	time.Sleep(time.Second)
	fmt.Println("broadcast...")
	cond.L.Lock()
	count -= 1
	cond.Broadcast()
	cond.L.Unlock()

	for i := 0; i < 5; i++ {
		<-ch
	}
}
