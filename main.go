package main

import (
	"fmt"
)

func main() {
	var a int
	fmt.Scanln(&a)
	for i := 1; i <= a; i++ {
		if i%3 == 0 {
			fmt.Println(i)
		}
	}
}
