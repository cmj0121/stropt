package stropt

import (
	"os"
	"testing"
)

type Foo struct {
	// the hidden field
	Ignore bool `-` //nolint
	// the innter field and should not process
	ignore bool //nolint

	// the flip option, store true/false
	Flip  bool
	Flip2 bool `shortcut:"f" name:"flip-2"`

	// the argument
	Age    uint `shortcut:"a"`
	Number int
	Name   string
	Price  float64   `shortcut:"p"`
	Point  complex64 `shortcut:"P"`
}

func Example() {
	foo := Foo{}
	parser := MustNew(&foo)

	parser.Usage(os.Stdout)
	// Output:
	// usage: foo [OPTION]
	//
	// options:
	//         --flip
	//      -f --flip-2
	//      -a --age
	//         --number
	//         --name
	//      -p --price
	//      -P --point
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
