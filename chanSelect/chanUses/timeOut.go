package main

import (
	"errors"
	"time"
)

func requestWithTimeout(timeout time.Duration) (int, error) {
	c := make(chan int)
	go doRequest(c)
	select {
	case data := <-c:
		return data, nil
	case <-time.After(timeout):
		return 0, errors.New("超时")
	}
}

func main() {

}
