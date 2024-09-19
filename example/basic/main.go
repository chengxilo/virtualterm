package main

import (
	"fmt"
	"github.com/chengxilo/virtualterm"
)

func main() {
	str := "hello\rvirtuaa\bl-terminal"
	vt := virtualterm.NewDefault()
	vt.Write([]byte(str))
	fmt.Println(str == "virtual-terminal")
	fmt.Println(vt.String() == "virtual-terminal")
	// Output:
	// false
	// true
}
