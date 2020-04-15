package strategy

import "fmt"

type Division struct {
}
type Addition struct {
}

type Multiplication struct {
}

type Subtraction struct {
}

func (Division) Compute(num1, num2 int) int {
	defer func() {
		if f := recover(); f != nil {
			fmt.Println(f)
			return
		}
	}()

	if num2 == 0 {
		panic("num2 must not be 0!")
	}
	return num1 / num2
}

func (p Subtraction) Compute(num1, num2 int) int {
	defer func() {
		if f := recover(); f != nil {
			fmt.Println(f)
			return
		}
	}()

	return num1 - num2
}

func (p Multiplication) Compute(num1, num2 int) int {
	defer func() {
		if f := recover(); f != nil {
			fmt.Println(f)
			return
		}
	}()

	return num1 * num2
}

func (p Addition) Compute(num1, num2 int) int {
	defer func() {
		if f := recover(); f != nil {
			fmt.Println(f)
			return
		}
	}()

	return num1 + num2
}

func NewStrategy(t string) (res Strategier) {
	switch t {
	case "s": // 减法
		res = Subtraction{}
	case "m": // 乘法
		res = Multiplication{}
	case "d": // 除法
		res = Division{}
	case "a": // 加法
		fallthrough
	default:
		res = Addition{}
	}

	return
}
