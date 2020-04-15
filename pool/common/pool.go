package common

import (
	"errors"
	"io"
	"log"
	"sync"
)

//一个安全的资源池，被管理的资源必须都实现io.Close接口
type Pool struct {
	//m是一个互斥锁，这主要是用来保证在多个goroutine访问资源时，池内的值是安全的。
	m sync.Mutex
	// res字段是一个有缓冲的通道，用来保存共享的资源，这个通道的大小，
	// 在初始化Pool的时候就指定的。注意这个通道的类型是io.Closer接口，
	// 所以实现了这个io.Closer接口的类型都可以作为资源，交给我们的资源池管理。
	res chan io.Closer
	// 一个函数类型，它的作用就是当需要一个新的资源时，可以通过这个函数创建，
	// 也就是说它是生成新资源的，至于如何生成、生成什么资源，是由使用者决定的，
	// 所以这也是这个资源池灵活的设计的地方
	factory func() (io.Closer, error)
	//表示资源池是否被关闭，如果被关闭的话，再访问是会有错误的
	closed bool
}

var ErrPoolClosed = errors.New("资源池已经关闭")

//创建一个资源池
func NewPool(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size的值太小了。")
	}
	return &Pool{
		factory: fn,
		res:     make(chan io.Closer, size),
	}, nil
}

//从资源池里获取一个资源
//Acquire方法可以从资源池获取资源，如果没有资源，则调用factory方法生成一个并返回。
//这里同样使用了select的多路复用，因为这个函数不能阻塞，可以获取到就获取，不能就生成一个。
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.res:
		log.Println("Acquire:共享资源")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("Acquire:新生成资源")
		return p.factory()
	}
}

//关闭资源池，释放资源
func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	//关闭通道，不让写入了
	p.closed = true

	//关闭通道里的资源
	close(p.res)
	for r := range p.res {
		_ = r.Close()
	}
}

func (p *Pool) Release(r io.Closer) {
	//保证该操作和Close方法的操作是安全的
	p.m.Lock()
	defer p.m.Unlock()
	//资源池都关闭了，就省这一个没有释放的资源了，释放即可
	if p.closed {
		_ = r.Close()
		return
	}

	select {
	case p.res <- r:
		log.Println("资源释放到池子了")
	default:
		log.Println("资源池满了，释放这个资源吧")
		_ = r.Close()
	}
}
