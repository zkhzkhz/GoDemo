package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Go语言在设计的时候，为了编写方便、效率高以及降低复杂度，被设计成为一门强类型的静态语言。
// 强类型意味着一旦定义了，它的类型就不能改变了；静态意味着类型检查在运行前就做了。
//
// 同时为了安全的考虑，Go语言是不允许两个指针类型进行转换的
func main() {
	i := 10
	ip := &i

	fmt.Println(i)
	fmt.Println(ip)
	fmt.Println(&i)
	fmt.Println(&ip)

	fmt.Println(reflect.TypeOf(ip))
	//cannot convert ip (type *int) to type *float64
	//var fp *float64 = (*float64)(ip)

	//unsafe.Pointer
	var fp = (*float64)(unsafe.Pointer(ip))
	*fp = *fp * 3
	fmt.Println(i)

	//unsafe.Pointer的4个规则。
	//
	//任何指针都可以转换为unsafe.Pointer
	//unsafe.Pointer可以转换为任何指针
	//uintptr可以转换为unsafe.Pointer
	//unsafe.Pointer可以转换为uintptr

	u := new(user)
	fmt.Println(*u)
	pName := (*string)(unsafe.Pointer(u))
	*pName = "张三"
	fmt.Println(unsafe.Offsetof(u.age))
	pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + unsafe.Offsetof(u.age)))
	*pAge = 20
	fmt.Println(*u)

	//这里我们可以看到，我们第二个偏移的表达式非常长，但是也千万不要把他们分段，不能像下面这样。
	//
	//	temp:=uintptr(unsafe.Pointer(u))+unsafe.Offsetof(u.age)
	//	pAge:=(*int)(unsafe.Pointer(temp))
	//	*pAge = 20
	//逻辑上看，以上代码不会有什么问题，但是这里会牵涉到GC，如果我们的这些临时变量被GC，那么导致的内存操作就错了，
	//我们最终操作的，就不知道是哪块内存了，会引起莫名其妙的问题。
}

type user struct {
	name string
	age  int
}

//unsafe是不安全的，所以我们应该尽可能少的使用它，比如内存的操纵，这是绕过Go本身设计的安全机制的，
//不当的操作，可能会破坏一块内存，而且这种问题非常不好定位。
//
//当然必须的时候我们可以使用它，比如底层类型相同的数组之间的转换；比如使用sync/atomic包中的一些函数时；
//还有访问Struct的私有字段时；该用还是要用，不过一定要慎之又慎。
//
//还有，整个unsafe包都是用于Go编译器的，不用运行时，在我们编译的时候，Go编译器已经把他们都处理了。
