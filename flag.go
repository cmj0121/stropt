package stropt

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

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
	if len(args) == 0 {
		err = fmt.Errorf("should pass %v", flag.Hint())
		return
	}

	switch kind := flag.Value.Type().Kind(); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var v int

		if v, err = strconv.Atoi(args[0]); err != nil {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
			return
		}

		flag.Value.SetInt(int64(v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var v int

		if v, err = strconv.Atoi(args[0]); err != nil {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
			return
		} else if v < 0 {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
			return
		}

		flag.Value.SetUint(uint64(v))
	case reflect.Float32, reflect.Float64:
		rat := &big.Rat{}
		if _, ok := rat.SetString(args[0]); !ok {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
			return
		}

		float, exact := rat.Float64()
		flag.Infof("convert %v to %v (exact: %v)", args[0], float, exact)
		flag.Value.SetFloat(float)
	case reflect.Complex64, reflect.Complex128:
		var cplx complex128

		if cplx, err = strconv.ParseComplex(args[0], 128); err != nil {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
			return
		}
		flag.Value.SetComplex(cplx)
	case reflect.String:
		// copy the argument to string
		flag.Value.SetString(args[0])
	default:
		err = fmt.Errorf("not support parse %v: %v", kind, flag.Value)
		return
	}

	n++
	return
}

// return the Tag of the field
func (flag *Flag) GetTag() (tag reflect.StructTag) {
	tag = flag.StructField.Tag
	return
}

// return the original name of the field
func (flag *Flag) GetName() (name string) {
	name = strings.ToLower(flag.StructField.Name)

	if value, ok := flag.StructField.Tag.Lookup(KEY_NAME); ok {
		// override the field's name
		name = strings.ToLower(value)
	}

	return
}

// return the shortcut of the field
func (flag *Flag) GetShortcut() (shortcut string) {
	if value, ok := flag.StructField.Tag.Lookup(KEY_SHORTCUT); ok {
		// override the field' shortcut, should be rune
		runes := []rune(value)
		switch {
		case len(runes) == 0:
			// no changed
		case len(runes) > 1:
			flag.Warnf("shortcut too large: %v (should be one and only one rune", value)
		default:
			shortcut = string(runes[0])
		}
	}

	return
}

//show the type hint of the field
func (flag *Flag) Hint() (hint string) {
	switch kind := flag.StructField.Type.Kind(); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		hint = "INT"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		hint = "UINT"
	case reflect.Float32, reflect.Float64:
		hint = "RAT"
	case reflect.Complex64, reflect.Complex128:
		hint = "COMPLEX"
	case reflect.String:
		hint = "STR"
	default:
		hint = "ARGS"
	}

	return
}

// the default value
func (flag *Flag) Default() (_default string) {
	if !flag.Value.IsZero() {
		// set the default value
		_default = fmt.Sprintf("%v", flag.Value)
	}

	return
}
