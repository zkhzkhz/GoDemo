package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	const N = 5
	var values [N]int32
	password := "test12345678"
	var wgA, wgB sync.WaitGroup
	wgA.Add(N)
	wgB.Add(1)

	for i := 0; i < N; i++ {
		i := i
		go func() {
			wgB.Wait() // 等待广播通知
			log.Printf("values[%v]=%v \n", i, values[i])
			wgA.Done()
		}()
	}

	// 下面这个循环保证将在上面的任何一个
	// wg.Wait调用结束之前执行。
	for i := 0; i < N; i++ {
		values[i] = 50 + rand.Int31n(50)
	}
	wgB.Done() // 发出一个广播通知
	wgA.Wait()
}
