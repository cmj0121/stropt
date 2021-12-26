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

	// the field of the struct
	reflect.StructField
}

func NewFlip(tracer *trace.Tracer, value reflect.Value, typ reflect.StructField) (flip *Flip, err error) {
	if value.Kind() != reflect.Bool {
		err = fmt.Errorf("cannot set %v as flip", value)
		return
	}

	flip = &Flip{
		Tracer:      tracer,
		Value:       value,
		StructField: typ,
	}

	return
}

// parse the pass argument. There is no argument would be used in Flip
func (flip *Flip) Parse(args ...string) (n int, err error) {
	flip.Tracef("parse no arguments. just flip the current value: %v", flip.Value)
	flip.Value.SetBool(!flip.Value.Bool())
	return
}

// return the Tag of the field
func (flip *Flip) GetTag() (tag reflect.StructTag) {
	tag = flip.StructField.Tag
	return
}

// return the original name of the field
func (flip *Flip) GetName() (name string) {
	name = flip.StructField.Name
	return
}
