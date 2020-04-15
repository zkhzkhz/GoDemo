package main

import (
	"log"
	"math/rand"
	"time"
)

type Customer struct {
	id int
}

type Bar1 chan Customer

func (bar Bar1) ServeCustomer1(c Customer) {
	log.Print("++ 顾客#", c.id, "开始饮酒", "当前人数", len(bar))
	time.Sleep(time.Second * time.Duration(3+rand.Intn(16)))
	log.Print("-- 顾客#", c.id, "离开酒吧")
	<-bar // 离开酒吧，腾出位子
}

func main() {
	rand.Seed(time.Now().UnixNano())

	bar24x7 := make(Bar1, 10)
	for customerId := 0; ; customerId++ {
		time.Sleep(time.Second * 2)
		customer := Customer{customerId}
		bar24x7 <- customer
		go bar24x7.ServeCustomer1(customer)
	}
	for {
		time.Sleep(time.Second)
	}
}
