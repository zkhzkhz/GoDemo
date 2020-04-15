package main

import (
	"fmt"
	"time"
)

//
//func main() {
//	value := make([]byte, 32*1024*1024)
//	if _, err := rand.Read(value); err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	done := make(chan struct{})
//
//	go func() {
//		sort.Slice(value, func(i, j int) bool {
//			return value[i] < value[j]
//		})
//		done <- struct{}{}
//	}()
//
//	<-done
//	fmt.Println(value[0], value[len(value)-1])
//}

func main() {
	done := make(chan struct{})
	go func() {
		fmt.Println("hello")
		time.Sleep(
			time.Second * 2)
		<-done
	}()
	done <- struct{}{}
	fmt.Println("Hello World")
}
