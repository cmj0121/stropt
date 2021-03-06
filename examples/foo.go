package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cmj0121/stropt"
)

type Sub struct {
	stropt.Model

	Age uint `shortcut:"a" desc:"store unsigned integer"`
}

func (sub Sub) Version_(stropt *stropt.StrOpt, _field stropt.Field) (err error) {
	os.Stdout.WriteString("sub (demo)")
	os.Exit(0)
	return
}

// the example struct and ready to setup by stropt
type Foo struct {
	stropt.Model

	Flip bool `shortcut:"f" desc:"store true/false value"`

	Number int  `shortcut:"n" attr:"required" desc:"store integer"`
	Age    uint `shortcut:"a" default:"21" desc:"store unsigned integer"`

	Name *string `desc:"save name as argument"`

	*Sub `json:"sub-1" desc:"sub-command"`
}

func main() {
	name := "defualt-name"

	foo := Foo{
		Name: &name,

		Sub: &Sub{
			Age: 123,
		},
	}

	parser := stropt.MustNew(&foo)
	parser.Version("foo (demo)")

	parser.Run()
	text, err := json.MarshalIndent(foo, "", "    ")
	fmt.Printf("error: %v\n%v\n", err, string(text))
	fmt.Printf("\n%#v\n", foo)
}
