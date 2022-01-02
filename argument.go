package stropt

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cmj0121/trace"
)

type Argument struct {
	// same as the Argument, but pass the shadow value to Argument
	*Flag

	reflect.Value

	// the shadow of the value, create and copy to original value
	shadow reflect.Value
}

func NewArgument(tracer *trace.Tracer, value reflect.Value, typ reflect.StructField) (arg *Argument, err error) {
	if !(value.Kind() == reflect.Ptr && typ.Type.Kind() == reflect.Ptr) {
		err = fmt.Errorf("%T cannot be the args: %v", value.Interface(), value.Kind())
		return
	}

	shadow := reflect.New(typ.Type.Elem())
	arg = &Argument{
		Value:  value,
		shadow: shadow,
	}
	arg.Flag, err = NewFlag(tracer, shadow.Elem(), typ)

	return
}

// parse the pass argument, should consumed one and only one argument
func (arg *Argument) Parse(args ...string) (n int, err error) {
	if n, err = arg.Flag.Parse(args...); err == nil {
		// copy the shadow to current value
		arg.Value.Set(arg.shadow)
	}

	return
}

// return the original name of the field
func (arg *Argument) GetName() (name string) {
	name = strings.ToLower(arg.StructField.Name)

	if value, ok := arg.StructField.Tag.Lookup(KEY_NAME); ok {
		// override the field's name
		name = strings.ToLower(value)
	}

	return
}

// arguments does not has shortcut, just return
func (arg *Argument) GetShortcut() (shortcut string) {
	return
}

// the default value
func (arg *Argument) Default() (_default string) {
	if !arg.Value.IsZero() {
		// set the default value
		_default = fmt.Sprintf("%v", arg.Value.Elem())
	}

	return
}
