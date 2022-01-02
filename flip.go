package stropt

import (
	"fmt"
	"reflect"
	"strings"

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

// return the name of the field
func (flip *Flip) GetName() (name string) {
	name = strings.ToLower(flip.StructField.Name)

	if value, ok := flip.StructField.Tag.Lookup(KEY_NAME); ok {
		// override the field's name
		name = strings.ToLower(value)
	}

	return
}

// return the shortcut of the field
func (flip *Flip) GetShortcut() (shortcut string) {
	if value, ok := flip.StructField.Tag.Lookup(KEY_SHORTCUT); ok {
		// override the field' shortcut, should be rune
		runes := []rune(value)
		switch {
		case len(runes) == 0:
			// no changed
		case len(runes) > 1:
			flip.Warnf("shortcut too large: %v (should be one and only one rune", value)
		default:
			shortcut = string(runes[0])
		}
	}

	return
}

// the default value
func (flip *Flip) Default() (_default string) {
	return
}
