package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
)

//多返回值 的 future/promise
func RandomGenerator() <-chan uint64 {
	c := make(chan uint64)
	go func() {
		rnds := make([]byte, 8)
		for {
			_, err := rand.Read(rnds)
			if err != nil {
				close(c)
			}
			c <- binary.BigEndian.Uint64(rnds)
		}
	}()
	return c
}

//数据聚合
func Aggregator(inputs ...<-chan uint64) <-chan uint64 {
	out := make(chan uint64)
	var wg sync.WaitGroup
	for _, input := range inputs {
		//input := input
		//go func() {
		//	for {
		//		out <- <-input
		//	}
		//}()

		wg.Add(1)
		in := input // 此行是必要的
		go func() {
			for {
				x, ok := <-in
				if ok {
					out <- x
				} else {
					wg.Done()
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out

	////要聚合的数据较少的情况下
	//out := make(chan uint64)
	//go func() {
	//	inA, inB := inputs[0], inputs[1]
	//	select {
	//	case v := <-inA:
	//		out <- v
	//	case v := <-inB:
	//		out <- v
	//	}
	//}()
}

//数据分流
func Divisor(input <-chan uint64, outputs ...chan<- uint64) {
	for _, out := range outputs {
		out := out
		go func() {
			out <- <-input
		}()
	}
}

//数据合成
func Composor(inA, inB <-chan uint64) <-chan uint64 {
	output := make(chan uint64)
	go func() {
		for {
			a1, b, a2 := <-inA, <-inB, <-inA
			output <- a1 ^ b&a2
		}
	}()
	return output
}

//数据复制/增殖
func Duplicator(in <-chan uint64) (<-chan uint64, <-chan uint64) {
	outA, outB := make(chan uint64), make(chan uint64)
	go func() {
		for {
			x := <-in
			outA <- x
			outB <- x
		}
	}()
	return outA, outB
}

//数据计算/分析
func Calculator(in <-chan uint64, out chan uint64) <-chan uint64 {
	if out == nil {
		out = make(chan uint64)
	}
	go func() {
		for {
			x := <-in
			out <- ^x
		}
	}()
	return out
}

//数据验证/过滤
func Filter(input <-chan uint64, output chan uint64) <-chan uint64 {
	if output == nil {
		output = make(chan uint64)
	}
	go func() {
		bigInt := big.NewInt(0)
		for {
			x := <-input
			bigInt.SetUint64(x)
			if bigInt.ProbablyPrime(1) {
				output <- x
			}
		}
	}()
	return output
}

//数据服务/存盘
func Printer(input <-chan uint64) {
	for {
		x, ok := <-input
		if ok {
			fmt.Println(x)
		} else {
			return
		}
	}
}

//组装数据流系统
func main() {
	Printer(Filter(
		Calculator(
			RandomGenerator(),
			nil),
		nil))
}
