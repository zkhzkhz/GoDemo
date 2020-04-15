package main

import (
	"fmt"
	"math/rand"
	"time"
)

func source1(c chan<- int32) {
	ra, rb := rand.Int31(), rand.Intn(3)+1

	time.Sleep(time.Duration(rb) * time.Second)
	select {
	case c <- ra:
	default:
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	c := make(chan int32, 1)
	for i := 0; i < 5; i++ {
		go source1(c)
	}
	rnd := <-c
	fmt.Println(rnd)
}
