package stropt

import (
	"os"
	"testing"
)

type Foo struct {
}

func Example() {
	foo := Foo{}
	parser := MustNew(&foo)

	parser.Usage(os.Stdout)
	// Output:
	// usage: foo
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