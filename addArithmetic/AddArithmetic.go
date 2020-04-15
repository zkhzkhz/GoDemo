package main

import "fmt"

func TwoSum1(nums []int, target int) []int {
	n := len(nums)

	for i, v := range nums {
		for j := i + 1; j < n; j++ {
			if v+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

func TwoNum2(nums []int, target int) []int {
	m := make(map[int]int, len(nums))

	for i, v := range nums {
		sub := target - v
		if j, ok := m[sub]; ok {
			return []int{i, j}
		} else {
			m[v] = i
		}
	}
	return nil
}

func main() {
	r1 := TwoSum1([]int{2, 7, 11, 15}, 9)
	fmt.Println(r1)
	r2 := TwoNum2([]int{11, 15, 2, 7}, 9)
	fmt.Println(r2)
}
