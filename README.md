# virtual-term

VirtualTerm is created to simulate a terminal,handle the special character such as '\r','\b'.
Write a string into VirtualTerm, you will know what would your string be like if you output it to stdout.

## Getting Started

example
```go
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
```

go test. This is why I want to create this repository.
If you don't use this,just use the str, all of them will fail.
```go
package test

import (
    "fmt"
    "github.com/chengxilo/virtualterm"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestVirtualTerm(t *testing.T) {
    str := "hello\rvirtuaa\bl-terminal"
    vt := virtualterm.NewDefault()
    vt.Write([]byte(str))
    assert.Equal(t, vt.String(), "virtual-terminal")
}

func ExampleVirtualTerm() {
    str := "hello\rvirtuaa\bl-terminal"
    vt := virtualterm.NewDefault()
    vt.Write([]byte(str))
    fmt.Print(vt.String())
    // Output:
    // virtual-terminal
}
```

## support
* '\r'
* '\n'
* '\b'



