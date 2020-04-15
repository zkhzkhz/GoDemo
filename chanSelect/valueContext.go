package main

import (
	"context"
	"fmt"
	"time"
)

var key string = "name"
//我们可以使用context.WithValue方法附加一对K-V的键值对，这里Key必须是等价性的，也就是具有可比性；Value值要是线程安全的。
//这样我们就生成了一个新的Context，这个新的Context带有这个键值对，在使用的时候，可以通过Value方法读取ctx.Value(key)。
//记住，使用WithValue传值，一般是必须的值，不要什么值都传递。
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	valueCtx := context.WithValue(ctx, key, "【监控1】")
	go watch1(valueCtx)

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()
	time.Sleep(5 * time.Second)
}

func watch1(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Value(key), "监控退出，停止le...")
			return
		default:
			fmt.Println(ctx.Value(key), "goroutine 监控中...")
			time.Sleep(2 * time.Second)
		}
	}
}
