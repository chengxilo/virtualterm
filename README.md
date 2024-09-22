# virtual-term

 VirtualTerm is created to simulate a terminal,handle the special character such as '\r','\b'.
Write a string into VirtualTerm, you will know what would your string be like if you output it to stdout.


## Now Support 😎

### Character Encoding
* ASCII
* UTF-8

### Control Characters
* `\b` backspace
* `\r` carriage return
* `\n` feed line
* `ESC[#A`	moves cursor up # lines
* `ESC[#B`	moves cursor down # lines
* `ESC[#C`	moves cursor right # columns
* `ESC[#D`	moves cursor left # columns
* `ESC[H`  moves cursor to home position

***WARNING:*** if you try to write not supported ESC to it, the output may can not be predicted

## install🛸
use go get to
```shell
go get -u github.com/chengxilo/virtualterm
```

## Getting Started 🤔

example
```go
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
    str,err := vt.String()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(str == "virtual-terminal")
    // Output:
    // false
    // true
}

```

Use `virtualterm.Process` function. You will not need to create a virtual terminal and input on your own.
```golang
package main

import (
    "fmt"
    "github.com/chengxilo/virtualterm"
)

func main() {
    str := "hello\rvirtuaa\bl-terminal"
    newS,_ := virtualterm.Process(str)
    fmt.Println(str == "virtual-terminal")
    fmt.Println(newS == "virtual-terminal")
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
    ns,_ := virtualterm.Process(str)
    assert.Equal(t, ns, "virtual-terminal")
}

func ExampleVirtualTerm() {
    str := "hello\rvirtuaa\bl-terminal"
    ns,_ := virtualterm.Process(str)
    fmt.Print(ns)
    // Output:
    // virtual-terminal
}
```

# Contribution 🎉

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements
- or whatever you want...