package singleton

import (
	"fmt"
	"sync"
)

var m *Manager
var lock *sync.Mutex = &sync.Mutex{}

func GetInstance() *Manager {
	// 加锁的代价是很大的,使用双重锁机制来提高效率
	// 代码只是稍作修改而已，不过我们用了两个判断，
	// 而且我们将同步锁放在了条件判断之后，这样做就避免了每次调用都加锁，
	// 提高了代码的执行效率。
	if m == nil {
		lock.Lock()
		defer lock.Unlock()
		if m == nil {
			m = &Manager{}
		}
	}

	return m
}

type Manager struct {
}

func (p Manager) Manage() {
	fmt.Println("manage...")
}
