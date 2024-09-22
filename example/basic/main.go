package main

import (
	"fmt"
	"github.com/chengxilo/virtualterm"
	"log"
)

func main() {
	str := "hello\rvirtuaa\bl-terminal"
	vt := virtualterm.NewDefault()
	vt.Write([]byte(str))
	fmt.Println(str == "virtual-terminal")
	str, err := vt.String()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str == "virtual-terminal")
	// Output:
	// false
	// true
}
