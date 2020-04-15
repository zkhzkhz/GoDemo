package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

//func Test_main(t *testing.T) {
//	var wg sync.WaitGroup
//	wg.Add(2)
//
//	go func() {
//		defer wg.Done()
//		for i := 1; i < 100; i++ {
//			fmt.Println("A:", i)
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		for i := 1; i < 100; i++ {
//			fmt.Println("B:", i)
//		}
//	}()
//	wg.Wait()
//}
var tcount int32
var twg2 sync.WaitGroup

func Test_madin(t *testing.T) {
	//for i := 1; i < 10; i++ {
	//	fmt.Println()
	//}
	//i := time.Now()
	//runtime.GOMAXPROCS(1)
	//var wg1 sync.WaitGroup
	//wg1.Add(2)
	//go func() {
	//	defer wg1.Done()
	//	for i := 1; i < 10000; i++ {
	//		fmt.Println("A:", i)
	//	}
	//}()
	//fmt.Println(runtime.NumGoroutine())
	//go func() {
	//	defer wg1.Done()
	//	for i := 1; i < 10000; i++ {
	//		fmt.Println("B:", i)
	//	}
	//}()
	//wg1.Wait()
	//beeLogger.Log.Info(string(i.Sub(time.Now())))

	twg2.Add(2)
	go incCountLock1()
	go incCountLock1()
	twg2.Wait()
	fmt.Println(tcount)
}

func incCountLock1() {
	defer twg2.Done()
	for i := 0; i < 2; i++ {
		//这里atomic.LoadInt32和atomic.StoreInt32两个函数，一个读取int32类型变量的值，
		// 一个是修改int32类型变量的值，
		// 这两个都是原子性的操作，Go已经帮助我们在底层使用加锁机制，
		// 保证了共享资源的同步和安全，所以我们可以得到正确的结果，
		// 这时候我们再使用资源竞争检测工具go build -race检查，也不会提示有问题了。
		//atomic包里还有很多原子化的函数可以保证并发下资源同步访问修改的问题，
		// 比如函数atomic.AddInt32可以直接对一个int32类型的变量进行修改，在原值的基础上再增加多少的功能，也是原子性的
		value := atomic.LoadInt32(&tcount)
		runtime.Gosched()
		value++
		atomic.StoreInt32(&tcount, value)
	}
}
