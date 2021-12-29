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

	// the shadow of the value, create and copy to original value
	shadow reflect.Value
}

func NewArgument(tracer *trace.Tracer, value reflect.Value, typ reflect.StructField) (args *Argument, err error) {
	if !(value.Kind() == reflect.Ptr && typ.Type.Kind() == reflect.Ptr) {
		err = fmt.Errorf("%T cannot be the args: %v", value.Interface(), value.Kind())
		return
	}

	shadow := reflect.New(typ.Type.Elem())
	args = &Argument{
		shadow: shadow,
	}
	args.Flag, err = NewFlag(tracer, shadow.Elem(), typ)

	return
}

// return the original name of the field
func (args *Argument) GetName() (name string) {
	name = strings.ToLower(args.StructField.Name)

	if value, ok := args.StructField.Tag.Lookup(KEY_NAME); ok {
		// override the field's name
		name = strings.ToLower(value)
	}

	return
}

// arguments does not has shortcut, just return
func (args *Argument) GetShortcut() (shortcut string) {
	return
}
