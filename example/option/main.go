package main

import (
	"fmt"
	"github.com/chengxilo/virtualterm"
)

func main() {
	//this will allow virtual term to output some information
	vt := virtualterm.NewOptions(virtualterm.OptionSilence(false))
	vt.Write([]byte("\033]a"))
	fmt.Println(vt.String())
}
