package main

import (
	"fmt"
	"os"
)

func main() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	go single(ch)
	single(ch)
}

func single(ch chan int) {
	for {
		i := <-ch
		fmt.Println(i)
		if i > 99 {
			os.Exit(0)
		}
		i++
		ch <- i
	}
}
