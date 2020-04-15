package singleton

import (
	"fmt"
	"sync"
)

var once sync.Once

// go中就中这么一种机制来保证代码只执行一次，而且不需要我们手工去加锁解锁。
// sync.Once，它有一个Do方法，在它中的函数go会只保证仅仅调用一次！
func GetInstanceGo() *Manager {
	once.Do(func() {
		m = &Manager{}
	})
	return m
}

func (p Manager) Manager() {
	fmt.Println("manage...")
}
