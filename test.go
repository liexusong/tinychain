package main

import "fmt"

func main() {
	test := make(chan struct{})
	nsent := 0
	for {
		select {
		case test <- struct{}{}:
			nsent++
		default:
			fmt.Println(nsent)
			return
		}
	}
}
