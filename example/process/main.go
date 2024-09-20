package main

import (
	"fmt"
	"github.com/chengxilo/virtualterm"
)

func main() {
	str := "hello\rvirtuaa\bl-terminal"
	newS, _ := virtualterm.Process(str)
	fmt.Println(str == "virtual-terminal")
	fmt.Println(newS == "virtual-terminal")
	// Output:
	// false
	// true
}
