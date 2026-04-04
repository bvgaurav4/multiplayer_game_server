package main

import (
	"fmt"
)

func lol(s []int) []int {
	s[1] = 69
	s = append(s, 69)
	return s
}
func main() {
	// s := [3]int{1, 2, 3}
	s := make([]int, 2, 5)
	s = append(s, 10)
	// s = append(s, 10)

	// s = append(s, 10)
	// s = append(s, 10)

	// s[4] = 10
	lol(s)
	fmt.Println(s)
	fmt.Println(cap(s), len(s))
}
