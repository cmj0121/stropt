package stropt

type Foo struct {
}

func Example() {
	foo := Foo{}
	MustNew(foo)
}
