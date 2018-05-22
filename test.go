package main

import "fmt"

func main() {
	test := make(chan interface{}, 1)
	close(test)

	select {
	case test <- "hello":
		fmt.Println(test)
	case <-test:
		fmt.Println("ok")

	default:
	}
}
