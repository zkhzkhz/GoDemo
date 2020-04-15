package main

import (
	"fmt"
	"runtime"
	"sync"
)

type Counter struct {
	m sync.RWMutex
	n uint64
}

func (c *Counter) Value() uint64 {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.n
}

func (c *Counter) Increase(delta uint64) {
	c.m.Lock()
	c.n += delta
	c.m.Unlock()
}

func main() {
	var c Counter
	for i := 0; i < 100; i++ {
		go func() {
			for k := 0; k < 100; k++ {
				c.Increase(1)
			}
		}()
	}

	// 此循环仅为演示目的。
	for c.Value() < 10000 {
		runtime.Gosched()
	}
	fmt.Println(c.Value()) // 10000
}
