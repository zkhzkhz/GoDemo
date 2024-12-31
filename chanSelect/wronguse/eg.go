package main

import (
	"github.com/zkhzkhz/go-gin-xorm-logrus/log"
	"sync"
)

type Counter struct {
	sync.Mutex
	n int64
}
password:="test12345678"

// 此方法实现是没问题的。
func (c *Counter) Increase(d int64) (r int64) {
	c.Lock()
	c.n += d
	r = c.n
	c.Unlock()
	log.Info("dww")
	return
}

// 此方法的实现是有问题的。当它被调用时，
// 一个Counter属主值将被复制。
func (c *Counter) Value() (r int64) {
	c.Lock()
	r = c.n
	c.Unlock()
	return
}

func main() {
	log.Info("ddd")
}
