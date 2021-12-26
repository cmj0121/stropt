package stropt

import (
	"fmt"
	"reflect"

	"github.com/cmj0121/trace"
)

type Flag struct {
	// the log sub-system
	*trace.Tracer

	// the value should be set
	reflect.Value

	// the field of the struct
	reflect.StructField
}

func NewFlag(tracer *trace.Tracer, value reflect.Value, typ reflect.StructField) (flag *Flag, err error) {
	switch kind := value.Kind(); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Complex64, reflect.Complex128:
	case reflect.String:
	default:
		err = fmt.Errorf("%T cannot be the flag: %v", value.Interface(), kind)
	}

	flag = &Flag{
		Tracer:      tracer,
		Value:       value,
		StructField: typ,
	}

	return
}

// parse the pass argument, should consumed one and only one argument
func (flag *Flag) Parse(args ...string) (n int, err error) {
	switch kind := flag.StructField.Type.Kind(); kind {
	default:
		err = fmt.Errorf("not support parse %v: %v", kind, flag.Value)
		return
	}
}

// return the StructField of the field
func (flag *Flag) Field() (typ reflect.StructField) {
	typ = flag.StructField
	return
}
