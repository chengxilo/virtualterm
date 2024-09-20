module github.com/chengxilo/virtualterm

go 1.22.3

require github.com/stretchr/testify v1.9.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// do not use them, they are not stable or loss of function
retract [v0.0.1,v1.0.1]