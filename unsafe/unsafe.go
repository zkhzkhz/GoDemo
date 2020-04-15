package main

import (
	"fmt"
	"unsafe"
)

// unsafe，顾名思义，是不安全的，Go定义这个包名也是这个意思，让我们尽可能的不要使用它，
// 如果你使用它，看到了这个名字，也会想到尽可能的不要使用它，或者更小心的使用它。
//
// 虽然这个包不安全，但是它也有它的优势，那就是可以绕过Go的内存安全机制，
// 直接对内存进行读写，所以有时候因为性能的需要，会冒一些风险使用该包，对内存进行操作。
func main() {
	//sizeof

	//对于和平台有关的int类型，这个要看平台是32位还是64位，会取最大的。比如我自己测试，
	// 以上输出，会发现int和int64的大小是一样的，因为我的是64位平台的电脑。
	fmt.Println(unsafe.Sizeof(true))
	fmt.Println(unsafe.Sizeof(int8(2)))
	fmt.Println(unsafe.Sizeof(int16(2)))
	fmt.Println(unsafe.Sizeof(int32(3233)))
	fmt.Println(unsafe.Sizeof(int64(24324)))
	fmt.Println(unsafe.Sizeof(int(324234234)))
	fmt.Println(unsafe.Sizeof(float32(42341.4123)))
	fmt.Println(unsafe.Sizeof(float64(42341.4123)))
	fmt.Println(unsafe.Sizeof(int(33)))
	fmt.Println(unsafe.Sizeof(map[string]int{"zhangsan": 22}))
	fmt.Println(unsafe.Sizeof(string("dee")))

	// allignof
	// Alignof返回一个类型的对齐值，也可以叫做对齐系数或者对齐倍数。
	// 对齐值是一个和内存对齐有关的值，合理的内存对齐可以提高内存读写的性能

	// 获取对齐值还可以使用反射包的函数，也就是说：unsafe.Alignof(x)等价于reflect.TypeOf(x).Align()。
	var b bool
	var i8 int8
	var i16 int16
	var i64 int64

	var f32 float32

	var s string

	var m map[string]string

	var p *int32
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(unsafe.Alignof(b))
	fmt.Println(unsafe.Alignof(i8))
	fmt.Println(unsafe.Alignof(i16))
	fmt.Println(unsafe.Alignof(i64))
	fmt.Println(unsafe.Alignof(f32))
	fmt.Println(unsafe.Alignof(s))
	fmt.Println(unsafe.Alignof(m))
	fmt.Println(unsafe.Alignof(p))

	// Offsetof函数
	// Offsetof函数只适用于struct结构体中的字段相对于结构体的内存位置偏移量。结构体的第一个字段的偏移量都是0.
	// 字段的偏移量，就是该字段在struct结构体内存布局中的起始位置(内存位置索引从0开始)。
	// 根据字段的偏移量，我们可以定位结构体的字段，进而可以读写该结构体的字段，哪怕他们是私有的，黑客的感觉有没有。
	// unsafe.Offsetof(u1.i)等价于reflect.TypeOf(u1).Field(i).Offset
	var u1 user1

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(unsafe.Offsetof(u1.b))
	fmt.Println(unsafe.Offsetof(u1.ss))
	fmt.Println(unsafe.Offsetof(u1.i))
	fmt.Println(unsafe.Offsetof(u1.j))

	//struct 大小
	var u2 user2
	var u3 user3
	var u4 user4
	var u5 user5
	var u6 user6
	//不同的字段顺序，最终决定struct的内存大小，所以有时候合理的字段顺序可以减少内存的开销。

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("u1 size is ", unsafe.Sizeof(u1))
	fmt.Println("u2 size is ", unsafe.Sizeof(u2))
	fmt.Println("u4 size is ", unsafe.Sizeof(u4))
	fmt.Println("u3 size is ", unsafe.Sizeof(u3))
	fmt.Println("u6 size is ", unsafe.Sizeof(u6))
	fmt.Println("u5 size is ", unsafe.Sizeof(u5))
}

type user1 struct {
	b  byte
	ss string
	i  int32
	j  int64
}
type user2 struct {
	b  byte
	j  int64
	i  int32
	ss string
}
type user3 struct {
	ss string
	b  byte
	i  int32
	j  int64
}
type user4 struct {
	j int64
	b byte
	i int32
}
type user5 struct {
	i int32
	j int64
	b byte
}
type user6 struct {
	i int32
	b byte
	j int64
}
