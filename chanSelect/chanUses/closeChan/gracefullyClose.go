package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const Max = 100000
	const NumReceivers = 100
	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	dataCh := make(chan int)

	//发送者
	go func() {
		for {
			if value := rand.Intn(Max); value == 4 {
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()

	//接收者
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			for value := range dataCh {
				log.Println(value)
			}
		}()
	}
	wgReceivers.Wait()
}
