package main

import "fmt"

// 斐波那契数列
func fibonacci() func() int64 {
	var x, y int64
	x, y = 0, 1
	return func() int64 {
		y, x = x+y, y
		return x
	}
}

func main() {
	f := fibonacci()
	for i := 90; i > 0; i-- {
		fmt.Printf("%20d ", f())
		if i%10 == 1 {
			fmt.Println()
		}
	}
}
