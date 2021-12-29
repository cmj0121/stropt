package main

import (
	"fmt"

	"github.com/cmj0121/stropt"
)

// the example struct and ready to setup by stropt
type Foo struct {
	stropt.Model

	Flip bool `shortcut:"f" desc:"store true/false value"`

	Number  int        `shortcut:"n" desc:"store integer"`
	Age     uint       `shortcut:"a" desc:"store unsigned integer"`
	Price   float64    `shortcut:"p" desc:"store float number, may rational number"`
	Point   complex128 `shortcut:"P" desc:"store the complex as point"`
	Message string     `shortcut:"m" desc:"store the raw string"`
}

func main() {
	foo := Foo{}
	parser := stropt.MustNew(&foo)
	parser.Version("foo (demo)")

	parser.Run()
	fmt.Printf("%#v\n", foo)
}
