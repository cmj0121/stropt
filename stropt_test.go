package stropt

import (
	"os"
	"testing"
)

type Inner struct {
	// the hidden field
	Ignore bool `-` //nolint
	// the innter field and should not process
	ignore bool `desc:"the ignore field"` //nolint

	InnerX int `name:"inner-x" desc:"the inner field"`
}

type Sub struct {
	// the hidden field
	Ignore bool `-` //nolint
	// the innter field and should not process
	ignore bool `desc:"the ignore field"` //nolint

	// the flip option, store true/false
	Flip bool `desc:"store true/false field"`
	// the argument
	Age uint `shortcut:"a"`
}

type Foo struct {
	Model

	// the hidden field
	Ignore bool `-` //nolint
	// the innter field and should not process
	ignore bool `desc:"the ignore field"` //nolint

	// the flip option, store true/false
	Flip  bool `desc:"store true/false field"`
	Flip2 bool `shortcut:"f" name:"flip-2"`

	// the argument
	Age    uint `shortcut:"a"`
	Number int  `desc:"store integer"`
	Name   string
	Price  float64   `shortcut:"p" desc:"store float number"`
	Point  complex64 `shortcut:"P"`

	// the embedded struct, should extend as the normal fields
	Inner

	// the sub-command
	*Sub `name:"subc"`
}

func Example() {
	foo := Foo{}
	parser := MustNew(&foo)

	parser.Usage(os.Stdout)
	// Output:
	// usage: foo [OPTION] [SUB-COMMAND]
	//
	// options:
	//      -h --help             show this help message and exit
	//      -v --version          show the version and exit
	//         --flip             store true/false field
	//      -f --flip-2
	//      -a --age
	//         --number           store integer
	//         --name
	//      -p --price            store float number
	//      -P --point
	//         --inner-x          the inner field
	//
	// sub-commands:
	//         subc
}

func TestInvalidType(t *testing.T) {
	cases := []interface{}{
		nil,
		true,
		false,
		1,
		1.2,
		'x',
		"test",
		rune(123),
	}

	for _, c := range cases {
		if _, err := New(c); err == nil {
			t.Errorf("expect cannot pass %T", c)
			continue
		}

		if _, err := New(&c); err == nil {
			t.Errorf("expect cannot pass %T", c)
			continue
		}
	}
}
