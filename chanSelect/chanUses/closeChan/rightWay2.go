package main

import "sync"

type MyChannel1 struct {
	C      chan struct{}
	closed bool
	mutex  sync.Mutex
}

func NewMyChannel1() *MyChannel1 {
	return &MyChannel1{C: make(chan struct{})}
}

func (mc *MyChannel1) SafeClose() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	if !mc.closed {
		close(mc.C)
		mc.closed = true
	}
}

func (mc *MyChannel1) IsClosed() bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	return mc.closed
}

func main() {

}
