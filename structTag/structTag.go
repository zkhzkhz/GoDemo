package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// 主要利用的就是User这个结构体对应的字段Tag，json解析的原理就是通过反射获得每个字段的tag，
// 然后把解析的json对应的值赋给他们。
// 利用字段Tag不光可以把Json字符串转为结构体对象，还可以把结构体对象转为Json字符串
func main() {
	var u User

	//反射获取字段Tag
	t := reflect.TypeOf(u)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fmt.Println(sf.Tag)
		fmt.Println(sf.Tag.Get("json"))
		fmt.Println(sf.Tag.Get("json"), ",", sf.Tag.Get("bson"))
	}
	h := `{"name":"dad","age":13}`
	err := json.Unmarshal([]byte(h), &u)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u)
	}
	newJson, err := json.Marshal(&u)
	fmt.Println(string(newJson))
}

type User struct {
	Name string `json:"name",bson:"b_name" `
	Age  int    `json:"age",bson:"b_age" `
}
