package main

import (
	"github.com/cmj0121/stropt"
)

// the example struct and ready to setup by stropt
type Foo struct {
	stropt.Model
}

func main() {
	foo := Foo{}
	parser := stropt.MustNew(&foo)
	parser.Version("foo (demo)")

	parser.Run()
}
