package main

import (
	"log"
	"time"
)

type T struct {
}

func worker(id int, ready <-chan T, done chan<- T) {
	<-ready
	log.Print("Worker#", id, "开始工作")
	time.Sleep(time.Second * time.Duration(id+1))
	log.Print("Worker#", id, "工作完成")
	done <- T{}
}

func main() {
	log.SetFlags(0)

	ready, done := make(chan T), make(chan T)
	go worker(0, ready, done)
	go worker(1, ready, done)
	go worker(2, ready, done)

	time.Sleep(time.Second * 3 / 2)

	close(ready)

	<-done
	<-done
	<-done
}
