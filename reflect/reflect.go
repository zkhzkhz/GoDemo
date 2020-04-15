package main

import (
	"fmt"
	"reflect"
)

// 在Go的反射定义中，任何接口都会由两部分组成的，一个是接口的具体类型，一个是具体类型对应的值。比如var i int = 3 ，
// 因为interface{}可以表示任何类型，所以变量i可以转为interface{}，所以可以把变量i当成一个接口，
// 那么这个变量在Go反射中的表示就是<Value,Type>，其中Value为变量的值3,Type变量的为类型int。
//
// 在Go反射中，标准库为我们提供两种类型来分别表示他们reflect.Value和reflect.Type，
// 并且提供了两个函数来获取任意对象的Value和Type。
func main() {
	u := User{"张三", 20}
	t := reflect.TypeOf(u)
	fmt.Println(t)
	v := reflect.ValueOf(u)
	fmt.Println(v)
	//因为在Go的反射中，把任意一个对象分为reflect.Value和reflect.Type，
	// 而reflect.Value又同时持有一个对象的reflect.Value和reflect.Type,
	// 所以我们可以通过reflect.Value的Interface方法实现还原。
	// 现在我们看看如何从一个reflect.Value获取对应的reflect.Type。
	user := v.Interface().(User)
	fmt.Println(user)
	t1 := v.Type()
	fmt.Println(t1)
	//获取底层类型
	fmt.Println(v.Kind())
	// 通过fmt.Printf函数为我们提供了简便的方法。
	fmt.Printf("%T\n", u)
	fmt.Printf("%v\n", u)

	//遍历字段和方法
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Name)
	}
	for i := 0; i < t.NumMethod(); i++ {
		fmt.Println(t.Method(i).Name)
	}

	//修改字段的值
	x := 2
	v1 := reflect.ValueOf(&x)
	v1.Elem().SetInt(100)
	fmt.Println(v1)
	fmt.Println(x)
	//修改结构体的值 CanSet方法可以帮助我们判断是否可以修改该对象。
	v2 := reflect.ValueOf(&u)
	fmt.Println(v2.Elem().CanSet())
	fmt.Println(v2)
	u1 := User{"zhaokh", 23}
	v2.Elem().Set(reflect.ValueOf(u1))
	fmt.Println(v2)
	fmt.Println(u)
	//动态调用方法 IsValid 来判断是否可用（存在）。
	v3 := reflect.ValueOf(u)
	mPrint := v3.MethodByName("Print")
	args := []reflect.Value{reflect.ValueOf("前缀")}
	fmt.Println(mPrint.Call(args))

}

type User struct {
	Name string
	Age  int
}

func (u User) Print(prefix string) {
	fmt.Printf("%s:Name is %s,Age is %d", prefix, u.Name, u.Age)
}
