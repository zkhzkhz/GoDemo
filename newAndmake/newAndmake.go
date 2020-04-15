package main

import (
	"fmt"
	"sync"
)

func main() {
	//var i *int
	//*i = 10
	//fmt.Println(*i)
	//它只接受一个参数，这个参数是一个类型，分配好内存后，返回一个指向该类型内存地址的指针。
	//同时请注意它同时把分配的内存置为零，也就是类型的零值。
	var i *int
	i = new(int)
	*i = 10
	fmt.Println(*i)

	u := new(user)
	u.lock.Lock()
	u.name = "zhangsan"
	u.lock.Unlock()
	fmt.Println(*u)
	//make也是用于内存分配的，但是和new不同，它只用于chan、map以及切片的内存创建，而且它返回的类型就是这三个类型本身，
	//而不是他们的指针类型，因为这三种类型就是引用类型，所以就没有必要返回他们的指针了。
	//
	//注意，因为这三种类型是引用类型，所以必须得初始化，但是不是置为零值，这个和new是不一样的。
	ints := make(map[string]int, 5)
	ints["fef"] = 3
	ints["fefs"] = 5
	fmt.Println(ints)

	//二者都是内存的分配（堆上），但是make只用于slice、map以及channel的初始化（非零值）；
	//而new用于类型的内存分配，并且内存置为零。所以在我们编写程序的时候，就可以根据自己的需要很好的选择了。

	//使用make的好处是可以指定len和cap，make(type,len,cap),合适的len和cap可以提升性能。

	//new这个内置函数，可以给我们分配一块内存让我们使用，但是现实的编码中，它是不常用的。
	//我们通常都是采用短语句声明以及结构体的字面量达到我们的目的，比如：
	//i1:=0
	//u1:=user{}
	//这样更简洁方便，而且不会涉及到指针这种比麻烦的操作。
	//make函数是无可替代的，我们在使用slice、map以及channel的时候，还是要使用make进行初始化，
	//然后才才可以对他们进行操作。
}

type user struct {
	lock sync.Mutex
	name string
	age  int
}
