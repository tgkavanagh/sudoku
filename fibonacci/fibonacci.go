package main

import (
	"fmt"
)

func fibrec(n int) uint64 {
	if n < 0 {
		return 0
	}

	if n < 2 {
		return 1
	}

	return fibrec(n-1) + fibrec(n-2)
}

func fib(n int) uint64 {
	if n < 0 {
		return 0
	}

	i := 0
	var sum uint64 = 1
	var prev uint64 = 1
	var prev1 uint64 = 0

	for i < n {
		sum = prev + prev1
		prev1 = prev
		prev = sum
		i++
	}

	return sum
}

func main() {
	i := 0

	for i < 100 {
		fmt.Printf("Fib(%d) = %d\n", i, fibrec(i))
		i++
	}
}
