package main

import (
	"log"
	"math/rand"
	"time"
)

type Customer1 struct {
	id int
}

type Bar2 chan Customer1

func (bar Bar2) ServeCustomer2(c Customer1) {
	log.Print("++ 顾客#", c.id, "开始饮酒", "当前人数", len(bar))
	time.Sleep(time.Second * time.Duration(30+rand.Intn(16)))
	log.Print("-- 顾客#", c.id, "离开酒吧")
	<-bar // 离开酒吧，腾出位子
}

func main() {
	rand.Seed(time.Now().UnixNano())

	bar24x7 := make(Bar2, 10)
	for customerId := 0; ; customerId++ {
		time.Sleep(time.Second * 1)
		customer := Customer1{customerId}
		select {
		case bar24x7 <- customer:
			go bar24x7.ServeCustomer2(customer)
		default:
			log.Print("顾客#", customerId, "不愿等待离去")
		}
	}
	for {
		time.Sleep(time.Second)
	}
}
