package stropt

import (
	"fmt"
	"reflect"

	"github.com/cmj0121/trace"
)

// the flip option, store true/false
type Flip struct {
	// the log sub-system
	*trace.Tracer

	// the value should be set
	reflect.Value

	// the description of field
	desc string
}

func NewFlip(tracer *trace.Tracer, value reflect.Value) (flip *Flip, err error) {
	if value.Kind() != reflect.Bool {
		err = fmt.Errorf("cannot set %v as flip", value)
		return
	}

	flip = &Flip{
		Tracer: tracer,
		Value:  value,
	}

	return
}

// parse the pass argument. There is no argument would be used in Flip
func (flip *Flip) Parse(args ...string) (n int, err error) {
	flip.Tracef("parse no arguments. just flip the current value: %v", flip.Value)
	flip.Value.SetBool(!flip.Value.Bool())
	return
}

// set the field description
func (flip *Flip) SetDescription(desc string) {
	flip.desc = desc
}

// get the field description
func (flip *Flip) Description() (desc string) {
	desc = flip.desc
	return
}
