package main

import "fmt"

type Noodles interface {
	Description() string
	Price() float32
}

type Egg struct {
	noodles Noodles
	name    string
	price   float32
}

func (p Egg) SetNoodles(noodles Noodles) {
	p.noodles = noodles
}

type Ramen struct {
	name  string
	price float32
}

func (p Ramen) Description() string {
	return p.name
}

func (p Ramen) Price() float32 {
	return p.price
}

func (p Egg) Description() string {
	return p.noodles.Description() + "+" + p.name
}

func (p Egg) Price() float32 {
	return p.noodles.Price() + p.price
}

type Sausage struct {
	noodles Noodles
	name    string
	price   float32
}

func (p Sausage) SetNoodles(noodles Noodles) {
	p.noodles = noodles
}

func (p Sausage) Description() string {
	return p.noodles.Description() + "+" + p.name
}

func (p Sausage) Price() float32 {
	return p.noodles.Price() + p.price
}

func main() {
	ramen := Ramen{name: "ramen", price: 10}
	fmt.Println(ramen.Description())
	fmt.Println(ramen.Price())

	egg := Egg{noodles: ramen, name: "egg", price: 2}
	fmt.Println(egg.Description())
	fmt.Println(egg.Price())

	sausage := Sausage{noodles: egg, name: "sausage", price: 2}
	fmt.Println(sausage.Description())
	fmt.Println(sausage.Price())

	egg2 := Egg{noodles: egg, name: "egg", price: 2}
	fmt.Println(egg2.Description())
	fmt.Println(egg2.Price())
}
