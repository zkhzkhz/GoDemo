package main

import (
	"log"
	"sync"
)

// 这个安全的Map被我们定义为一个SynchronizedMap的结构体，这个结构体里有两个字段，一个是读写锁rw,一个是存储数据的data，
// data是map类型。
// 然后就是给SynchronizedMap定义一些方法，如果这些方法是增删改的，就要使用写锁，如果是只读的，就使用读锁，
// 这样就保证了我们数据data在多个goroutine下的安全性。
// 有了这个安全的Map我们就可以在多goroutine下增删改查数据了，都是安全的。

//安全的Map
type SynchronizedMap struct {
	rw   *sync.RWMutex
	data map[interface{}]interface{}
}

//存储操作
func (sm *SynchronizedMap) Put(k, v interface{}) {
	sm.rw.Lock()
	defer sm.rw.Unlock()

	sm.data[k] = v
}

//获取操作
func (sm *SynchronizedMap) Get(k interface{}) interface{} {
	sm.rw.RLock()
	defer sm.rw.RUnlock()

	return sm.data[k]
}

//删除操作
func (sm *SynchronizedMap) Delete(k interface{}) {
	sm.rw.Lock()
	defer sm.rw.Unlock()

	delete(sm.data, k)
}

//遍历Map，并且把遍历的值给回调函数，可以让调用者控制做任何事情
func (sm *SynchronizedMap) Each(cb func(interface{}, interface{})) {
	sm.rw.RLock()
	defer sm.rw.RUnlock()

	for key, value := range sm.data {
		cb(key, value)
	}
}

//生成初始化一个SynchronizedMap
func NewSynchronizedMap() *SynchronizedMap {
	return &SynchronizedMap{
		rw:   new(sync.RWMutex),
		data: make(map[interface{}]interface{}),
	}
}

func init() {
	log.SetPrefix("【UserCenter】")
	log.SetFlags(log.Ltime | log.Ldate | log.Lshortfile)
}
func main() {
	//参数skip表示跳过栈帧数，0表示不跳过，也就是runtime.Caller的调用者。1的话就是再向上一层，表示调用者的调用者。
	//log日志包里使用的是2，也就是表示我们在源代码中调用log.Print、log.Fatal和log.Panic这些函数的调用者。
	//以main函数调用log.Println为例，是main->log.Println->*Logger.Output->runtime.Caller这么一个方法调用栈，
	// 所以这时候，skip的值分别代表：
	//
	//0表示*Logger.Output中调用runtime.Caller的源代码文件和行号
	//1表示log.Println中调用*Logger.Output的源代码文件和行号
	//2表示main中调用log.Println的源代码文件和行号
	//所以这也是log包里的这个skip的值为什么一直是2的原因。
	log.Panic("飞雪无情的博客:", "http://www.flysnow.org")
	log.Printf("飞雪无情的微信公众号：%s\n", "flysnow_org")
}
