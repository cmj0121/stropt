package main

import (
	"fmt"

	"github.com/cmj0121/stropt"
)

// the example struct and ready to setup by stropt
type Foo struct {
	stropt.Model

	Flip bool `shortcut:"f" desc:"store true/false value"`
}

func main() {
	foo := Foo{}
	parser := stropt.MustNew(&foo)
	parser.Version("foo (demo)")

	parser.Run()
	fmt.Printf("%#v\n", foo)
}
