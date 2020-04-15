package main

import (
	"fmt"
	"math/rand"
	"sync"
)

//读写锁

// 互斥锁的根本就是当一个goroutine访问的时候，其他goroutine都不能访问，这样肯定保证了资源的同步，
// 避免了竞争，不过也降低了性能

// 当我们读取一个数据的时候，如果这个数据永远不会被修改，那么其实是不存在资源竞争的问题的，
// 因为数据是不变的，不管怎么读取，多少goroutine同时读取，都是可以的。
// 所以其实读取并不是问题，问题主要是修改，修改的数据要同步，这样其他goroutine才可以感知到。
// 所以真正的互斥应该是读取和修改、修改和修改之间，读取和读取是没有互斥操作的。

var count int
var wg sync.WaitGroup

// 以上我们定义了一个共享的资源count，并且声明了2个函数进行读写read和write，在main函数的测试中，
// 我们同时启动了5个读写goroutine进行读写操作，通过打印的结果来看，写入操作是处于竞争状态的，
// 有的写入操作被覆盖了。通过go build -race也可以看到更明细的竞争态。
func main() {
	wg.Add(10)

	for i := 0; i < 5; i++ {
		//go read(i)
		go read1(i)
	}

	for i := 0; i < 5; i++ {
		//go write(i)
		go write1(i)
	}
	wg.Wait()
}

func read(n int) {
	fmt.Printf("读goroutine %d 正在读取...\n", n)

	v := count
	fmt.Printf("读goroutine %d 读取结束，值为：%d\n", n, v)
	wg.Done()
}

func write(n int) {
	fmt.Printf("写goroutine %d 正在写入...\n", n)
	v := rand.Intn(1000)

	count = v
	fmt.Printf("写goroutine %d写入结束，新值为：%d\n", n, v)
	wg.Done()
}

var rw sync.RWMutex

func read1(n int) {
	rw.RLock()
	fmt.Printf("读goroutine %d 正在读取...\n", n)

	v := count

	fmt.Printf("读goroutine %d 读取结束，值为：%d\n", n, v)
	wg.Done()
	rw.RUnlock()
}

func write1(n int) {
	rw.Lock()
	fmt.Printf("写goroutine %d 正在写入...\n", n)
	v := rand.Int()

	count = v

	fmt.Printf("写goroutine %d 写入结束，新值为：%d\n", n, v)
	wg.Done()
	rw.Unlock()
}
