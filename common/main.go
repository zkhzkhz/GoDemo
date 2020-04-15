package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	//go func() {
	//	for{}
	//}()
	//time.Sleep(time.Second)
	//fmt.Println("the answer to  life:",42)
	//go println("你好，并发")

	//for {}
	//select {}
	//<-make(chan bool)

	//runtime.Gosched()
	//for i := 0; i < 10000*100; i++ {
	//	go printSum(i)
	//}
	// main函数在退出前需要从done管道取一个消息，后台任务在将消息放入done管道前必须先完成自己的输出任务。因此，main
	// 函数成功取到消息时，后台的输出任务确定已经完成了，main函数也就可以放心退出了。
	done := make(chan bool)
	go func() {
		println("he")
		done <- true
	}()

	<-done
}
func println(s string) {
	print(s + "\n")
}

func printSum(n int) {
	fmt.Printf("sum(%[1]d): %[1]d\n", n)

}

func sum(n int) int {
	return sum(n-1) + n
}
