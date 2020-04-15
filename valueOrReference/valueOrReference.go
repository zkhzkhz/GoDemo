package main

import "fmt"

//值传递 传递的是值的拷贝
func main() {
	i := 10
	ip := &i
	fmt.Printf("%p", &i)
	fmt.Println()
	fmt.Printf("%p\n", &ip)
	fmt.Println(*ip)
	fmt.Println(**(&ip))
	modify(ip)
	fmt.Println(i)

	//Go语言通过make函数，字面量的包装，为我们省去了指针的操作，
	//让我们可以更容易的使用map。这里的map可以理解为引用类型，但是记住引用类型不是传引用
	//map
	persons := make(map[string]int)
	persons["张三"] = 19

	mp := &persons
	fmt.Printf("%p\n", mp)
	modifyMap(persons)
	fmt.Println(persons)

	//struct
	p := Person{"张三"}
	modifyStru(&p)
	fmt.Println(p)

	ages := []int{6, 6, 6}
	//通过源代码发现，对于chan、map、slice等被当成指针处理，通过value.Pointer()获取对应的值的指针
	fmt.Printf("%p\n", ages)
	fmt.Printf("%p\n", &ages)
	modifySli(ages)
	fmt.Println(ages)
}

func modifySli(ages []int) {
	fmt.Printf("%p\n", ages)
	ages[0] = 1
}

func modifyMap(p map[string]int) {
	fmt.Printf("%p\n", &p)
	p["张三"] = 20
}
func modify(ip *int) {
	fmt.Printf("%p", &ip)
	fmt.Println()
	*ip = 1
}

type Person struct {
	Name string
}

func modifyStru(p *Person) {
	p.Name = "里斯"
}

//最终我们可以确认的是Go语言中所有的传参都是值传递（传值），
//都是一个副本，一个拷贝。因为拷贝的内容有时候是非引用类型（int、string、struct等这些），
//这样就在函数中就无法修改原内容数据；有的是引用类型（指针、map、slice、chan等这些），这样就可以修改原内容数据。
//
//是否可以修改原内容数据，和传值、传引用没有必然的关系。在C++中，传引用肯定是可以修改原内容数据的，
//在Go语言里，虽然只有传值，但是我们也可以修改原内容数据，因为参数是引用类型。
