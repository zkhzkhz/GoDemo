package main

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

type person struct {
	Name string
}

func main() {
	i := 0
	s := "哈哈"
	spew.Dump(i, s)

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()

	m := map[int]string{1: "1", 2: "2"}
	e := errors.New("嘿嘿，错误")
	p := person{Name: "张三"}
	spew.Dump(i, s, m, e, p)
}
