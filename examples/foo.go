package main

import (
	"encoding/json"
	"fmt"

	"github.com/cmj0121/stropt"
)

// the example struct and ready to setup by stropt
type Foo struct {
	stropt.Model

	Flip bool `shortcut:"f" desc:"store true/false value"`

	Number  int     `shortcut:"n" desc:"store integer"`
	Age     uint    `shortcut:"a" desc:"store unsigned integer"`
	Price   float64 `shortcut:"p" desc:"store float number, may rational number"`
	Message string  `shortcut:"m" desc:"store the raw string"`

	Name   *string `desc:"save name as argument"`
	Amount *int    `desc:"save int as argument"`
}

func main() {
	foo := Foo{}
	parser := stropt.MustNew(&foo)
	parser.Version("foo (demo)")

	parser.Run()
	text, err := json.MarshalIndent(foo, "", "    ")
	fmt.Printf("error: %v\n%v\n", err, string(text))
}
