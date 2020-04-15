package main

import (
	"./strategy"
	"./compute"
	"flag"
	"fmt"
)

var stra *string = flag.String("type", "a", "input the strategy")
var num1 *int = flag.Int("num1", 1, "input num1")
var num2 *int = flag.Int("num2", 1, "input num2")

func init() {
	flag.Parse()
}

// 策略模式还是算比较容易理解的，策略模式的核心就是将容易变动的代码从主逻辑中分离出来，通过一个接口来规范它们的形式，
// 在主逻辑中将任务委托给策略。这样做既减少了我们对主逻辑代码修改的可能性，也增加了系统的可扩展性。一定要记得哦，
// 我们的代码要往对扩展开发，对修改关闭这条设计原则上努力！
func main() {
	com := compute.Computer{Num1: *num1, Num2: *num2}
	strate := strategy.NewStrategy(*stra)

	com.SetStrategy(strate)
	fmt.Println(com.Do())
}
