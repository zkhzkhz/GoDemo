package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const Max = 1000000
	const NumReceivers = 10
	const NumSenders = 1000
	const NumThirdParties = 15

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int)   // 仍就不会被关闭
	middleCh := make(chan int) // 将被关闭
	closing := make(chan string)
	closed := make(chan struct{})

	var stoppedBy string

	stop := func(by string) {
		select {
		case closing <- by:
			<-closed
		case <-closed:
		}
	}

	// 中间层
	go func() {
		exit := func(v int, needSend bool) {
			close(closed)
			if needSend {
				dataCh <- v
			}
			close(dataCh)
		}

		for {
			select {
			case stoppedBy = <-closing:
				exit(0, false)
				return
			case v := <-middleCh:
				select {
				case stoppedBy = <-closing:
					exit(v, true)
					return
				case dataCh <- v:
				}
			}
		}
	}()

	// 一些第三方协程
	for i := 0; i < NumThirdParties; i++ {
		go func(id string) {
			r := 1 + rand.Intn(3)
			time.Sleep(time.Duration(r) * time.Second)
			stop("3rd-party#" + id)
		}(strconv.Itoa(i))
	}

	// 发送者
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(Max)
				if value == 0 {
					stop("sender#" + id)
					return
				}

				select {
				case <-closed:
					return
				default:
				}

				select {
				case <-closed:
					return
				case middleCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// 接收者
	for range [NumReceivers]struct{}{} {
		go func() {
			defer wgReceivers.Done()

			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	// ...
	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
