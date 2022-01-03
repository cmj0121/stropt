# Struct Option #
[![GitHub Action][0]][1]

The struct-based command-line parser for Go-lang.

## Example ##
```go
package main

import (
	"github.com/cmj0121/stropt"
)

type Foo struct {
	stropt.Model

	Number  int     `shortcut:"n" desc:"store integer"`
	Age     uint    `shortcut:"a" default:"21" desc:"store unsigned integer"`
	Price   float64 `shortcut:"p" default:"12.34" desc:"store float number, may rational number"`
	Message string  `shortcut:"m" desc:"store the raw string"`
}

func main() {
	foo := Foo{}
	parser := stropt.MustNew(&foo)
	parser.Run()
}
```

[0]: https://github.com/cmj0121/stropt/actions/workflows/pipeline.yml/badge.svg
[1]: https://github.com/cmj0121/stropt/actions
