package main

import "fmt"

func main() {
	a := make(map[string][]string)
	a["hi"] = []string{"a", "b"}
	b := a["hi"]
	b = append(b, "c")
	a["hi"] = b
	fmt.Println(a)
}
