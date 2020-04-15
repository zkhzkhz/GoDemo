package main

import (
	"math/rand"
	"testing"
)

//func BenchmarkTwoSum1(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		TwoSum1([]int{2, 7, 11, 15}, 9)
//	}
//}
//func BenchmarkTwoSum2(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		TwoNum2([]int{2, 7, 11, 15}, 9)
//	}
//}
const N = 1000

func BenchmarkTwoSum1(b *testing.B) {
	var nums []int
	for i := 0; i < N; i++ {
		nums = append(nums, rand.Int())
	}
	nums = append(nums, 7, 2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TwoSum1(nums, 9)
	}
}

func BenchmarkTwoSum2(b *testing.B) {
	var nums []int
	for i := 0; i < N; i++ {
		nums = append(nums, rand.Int())
	}
	nums = append(nums, 7, 2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TwoNum2(nums, 9)
	}
}
