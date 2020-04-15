package main

import "go-gin-xorm-logrus/log"

func main() {
	var c = make(chan string, 1)
	c <- "da"
	b := c
	log.Info(b)
	log.Info(c)
	log.Info(&b)
	log.Info(&c)
	log.Info(cap(c))
	v := <-c
	log.Info(v)
	//v, sentBeforeClosed := <-c
	//log.Info(sentBeforeClosed)
	var ch chan string
	log.Info(len(ch))
	log.Info(cap(ch))
	ch = make(chan string)
	//ch <- "1"
	//<-ch
	log.Info(b == c, &b == &c)
}
