package main

import "fmt"

func main() {
	foo, bar := make(chan struct{}), make(chan struct{})
	close(foo)
	close(bar) // 仅为演示目的
	x, y := 0.0, 0.0
	f := func() { x++ }
	g := func() { y++ }
	for i := 0; i < 100000; i++ {
		select {
		case <-foo:
			f()
		case <-foo:
			f()
		case <-bar:
			g()
		}
	}
	fmt.Println(x / y) // 大致为2
}
