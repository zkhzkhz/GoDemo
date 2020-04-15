package main

import "fmt"

type Phone interface {
	ShowBrand()
}

type IPhone struct {
}

func (phone IPhone) ShowBrand() {
	fmt.Println("[Phone Brand]:Apple")
}

type HPhone struct {
}

func (phone HPhone) ShowBrand() {
	fmt.Println("[Phone Brand]:Huawei")
}

type XPhone struct {
}

func (phone XPhone) ShowBrand() {
	fmt.Println("[Phone Brand]:Xiaomi")
}

type PhoneFactory struct {
}

func (factory PhoneFactory) CreatePhone(brand string) Phone {
	switch brand {
	case "HW":
		return new(HPhone)
	case "XM":
		return new(XPhone)
	case "PG":
		return new(IPhone)
	default:
		return nil
	}
}

func main() {

	var phone Phone
	factory := new(PhoneFactory)

	phone = factory.CreatePhone("HW")
	phone.ShowBrand()

	phone = factory.CreatePhone("XM")
	phone.ShowBrand()

	phone = factory.CreatePhone("PG")
	phone.ShowBrand()
}
