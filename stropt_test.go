package stropt

import (
	"fmt"
	"os"
	"strconv"
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
	Age *uint `attr:"required"`
}

type Foo struct {
	LogModel

	// the hidden field
	Ignore bool `-` //nolint
	// the innter field and should not process
	ignore bool `desc:"the ignore field"` //nolint

	// the flip option, store true/false
	Flip  bool `desc:"store true/false field"`
	Flip2 bool `shortcut:"f" name:"flip-2"`

	// the argument
	Age    uint       `attr:"required" shortcut:"a" default:"21" desc:"age"`
	Number int        `desc:"store integer"`
	Name   string     `default:"mock-name" desc:"name"`
	Price  float64    `shortcut:"p" desc:"store float number"`
	Point  complex128 `shortcut:"P"`

	// the embedded struct, should extend as the normal fields
	Inner

	// arguments
	Message *string `desc:"store string as position field"`
	Amount  *int    `desc:"store int as position field"`

	// the sub-command
	*Sub `name:"subc"`
}

func Example() {
	foo := Foo{
		Price: 12.34,
	}
	parser := MustNew(&foo)

	parser.Usage(os.Stdout)
	// Output:
	// usage: foo [OPTION] [ARGS] ...
	//
	// options:
	//      -h --help             show this help message and exit
	//      -v --version          show the version and exit
	//      -l --level STR        the log level [error warn info debug trace]
	//         --flip             store true/false field
	//      -f --flip-2
	//      -a --age UINT         age [default: 21] (required)
	//         --number INT       store integer
	//         --name STR         name [default: mock-name]
	//      -p --price RAT        store float number [default: 12.34]
	//      -P --point CPLX
	//         --inner-x INT      the inner field
	//
	// arguments:
	//         message STR        store string as position field
	//         amount INT         store int as position field
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

func TestParseFlip(t *testing.T) {
	foo := Foo{}
	parser := MustNew(&foo)

	parser.Parse("--flip")            //nolint
	parser.Parse("--flip", "--flip")  //nolint
	parser.Parse("-f", "-ff", "-fff") //nolint
}

func TestParseFlag(t *testing.T) {
	foo := Foo{}
	parser := MustNew(&foo)

	parser.Parse("--age", "12") //nolint
	parser.Parse("-a", "12")    //nolint
}

func TestParseWithDoubleDash(t *testing.T) {
	foo := Foo{}
	parser := MustNew(&foo)

	parser.Parse("--flip", "--", "--flip")           //nolint
	parser.Parse("--flip", "--flip", "--", "--flip") //nolint
	parser.Parse("-f", "-ff", "--", "-fff")          //nolint
}

func TestParseInt(t *testing.T) {
	foo := &Foo{}
	parser := MustNew(foo)

	number := 123
	if _, err := parser.Parse("--number", strconv.Itoa(number)); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v (%v)", err, strconv.Itoa(number))
	} else if foo.Number != number {
		// parse int fail
		t.Errorf("parse int %v (%v): %v", number, strconv.Itoa(number), foo.Number)
	}

	number = -123
	if _, err := parser.Parse("--number", strconv.Itoa(number)); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v", err)
	} else if foo.Number != number {
		// parse int fail
		t.Errorf("parse int %v (%v): %v", number, strconv.Itoa(number), foo.Number)
	}

	invalid_number := "123a"
	if _, err := parser.Parse("--number", invalid_number); err == nil {
		// expect parse int fail
		t.Errorf("parse %v without error", invalid_number)
	}
}

func TestParseUint(t *testing.T) {
	foo := &Foo{}
	parser := MustNew(foo)

	age := 123
	if _, err := parser.Parse("--age", strconv.Itoa(age)); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v (%v)", err, strconv.Itoa(age))
	} else if foo.Age != uint(age) {
		// parse int fail
		t.Errorf("parse int %v (%v): %v", age, strconv.Itoa(age), foo.Age)
	}

	age = -123
	if _, err := parser.Parse("--age", strconv.Itoa(age)); err == nil {
		// expect parse int fail
		t.Errorf("parse %v without error", age)
	}

	invalid_age := "123a"
	if _, err := parser.Parse("--age", invalid_age); err == nil {
		// expect parse int fail
		t.Errorf("parse %v without error", invalid_age)
	}
}

func TestParseFloat(t *testing.T) {
	foo := &Foo{}
	parser := MustNew(foo)

	price := 123.456
	if _, err := parser.Parse("--price", fmt.Sprintf("%v", price)); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v (%v)", err, fmt.Sprintf("%v", price))
	} else if foo.Price != price {
		// parse int fail
		t.Errorf("parse int %v (%v): %v", price, fmt.Sprintf("%v", price), foo.Price)
	}

	price = -123.456
	if _, err := parser.Parse("--price", fmt.Sprintf("%v", price)); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v", err)
	} else if foo.Price != price {
		// parse int fail
		t.Errorf("parse int %v (%v): %v", price, fmt.Sprintf("%v", price), foo.Price)
	}

	price = 0.5
	if _, err := parser.Parse("--price", "1/2"); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v", err)
	} else if foo.Price != price {
		// parse int fail
		t.Errorf("parse int %v (1/2): %v", price, foo.Price)
	}

	invalid_price := "123a"
	if _, err := parser.Parse("--price", invalid_price); err == nil {
		// expect parse int fail
		t.Errorf("parse %v without error", invalid_price)
	}
}

func TestParseComplex(t *testing.T) {
	foo := &Foo{}
	parser := MustNew(foo)

	point := 1 + 2i
	if _, err := parser.Parse("--point", fmt.Sprintf("%v", point)); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v (%v)", err, fmt.Sprintf("%v", point))
	} else if foo.Point != point {
		// parse int fail
		t.Errorf("parse int %v (%v): %v", point, fmt.Sprintf("%v", point), foo.Point)
	}

	point = -2 - 5i
	if _, err := parser.Parse("--point", fmt.Sprintf("%v", point)); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v", err)
	} else if foo.Point != point {
		// parse int fail
		t.Errorf("parse int %v (%v): %v", point, fmt.Sprintf("%v", point), foo.Point)
	}

	invalid_point := "123a"
	if _, err := parser.Parse("--point", invalid_point); err == nil {
		// expect parse int fail
		t.Errorf("parse %v without error", invalid_point)
	}
}

func TestParseString(t *testing.T) {
	foo := &Foo{}
	parser := MustNew(foo)

	name := "message/訊息/メッセージ"
	if _, err := parser.Parse("--name", name); err != nil {
		// cannot parse int
		t.Errorf("cannot parse int: %v (%v)", err, name)
	} else if foo.Name != name {
		// parse int fail
		t.Errorf("parse int %v : %v", name, foo.Name)
	}
}

func TestParseArgument(t *testing.T) {
	foo := &Foo{}
	parser := MustNew(foo)

	if _, err := parser.Parse("message"); err != nil {
		// cannot parse arguments
		t.Errorf("cannot parse argument: %v", err)
	} else if foo.Message == nil || *foo.Message != "message" {
		// parse argument fail
		t.Errorf("parse message fail: %#v", foo.Message)
	}

	if _, err := parser.Parse("123"); err != nil {
		// cannot parse arguments
		t.Errorf("cannot parse argument: %v", err)
	} else if foo.Amount == nil || *foo.Amount != 123 {
		// parse argument fail
		t.Errorf("parse amount fail: %#v", foo.Amount)
	}

	if _, err := parser.Parse("ccc"); err == nil {
		// expect cannot parse argument
		t.Errorf("expect cannot parse extra argument")
	}
}

func TestRequired(t *testing.T) {
	foo := &Foo{}
	parser := MustNew(foo)

	if _, err := parser.Parse("subc"); err == nil {
		t.Fatalf("expect required is works")
	} else if _, err := parser.Parse("subc", "--flip"); err == nil {
		t.Fatalf("expect required is works")
	} else if _, err := parser.Parse("subc", "12"); err != nil {
		t.Errorf("cannot setup the required field: %v", err)
	}
}
