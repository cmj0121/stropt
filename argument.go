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

	_default string
}

func NewArgument(tracer *trace.Tracer, value reflect.Value, typ reflect.StructField) (arg *Argument, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		shadow := reflect.New(typ.Type.Elem())
		arg = &Argument{
			Value:  value,
			shadow: shadow,
		}
		arg.Flag, err = NewFlag(tracer, shadow.Elem(), typ)
	default:
		shadow := reflect.New(typ.Type).Elem()

		arg = &Argument{
			Value:  value,
			shadow: shadow,
		}
		arg.Flag, err = NewFlag(tracer, shadow, typ)
	}

	if v, ok := arg.Tag.Lookup(KEY_DEFAULT); ok {
		// set default if defined as tag
		arg._default = v
	} else if !value.IsZero() {
		// only set the default if value is not Zero
		arg._default = fmt.Sprintf("%v", value)
	}
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
	_default = arg._default
	return
}

// check the field set or not
func (arg *Argument) IsZero() bool {
	return arg.Value.IsZero()
}
