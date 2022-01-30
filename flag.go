package stropt

import (
	"fmt"
	"math/big"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cmj0121/trace"
)

type Flag struct {
	// the log sub-system
	*trace.Tracer

	// the value should be set
	reflect.Value

	// the field of the struct
	reflect.StructField

	// the possible and valid value can be set
	choise []string

	// the default value
	_default string
}

func NewFlag(tracer *trace.Tracer, value reflect.Value, typ reflect.StructField) (flag *Flag, err error) {
	flag = &Flag{
		Tracer:      tracer,
		Value:       value,
		StructField: typ,
	}

	if v, ok := flag.Tag.Lookup(KEY_DEFAULT); ok {
		// set default if defined as tag
		flag._default = v
	} else if !value.IsZero() {
		// only set the default if value is not Zero
		flag._default = fmt.Sprintf("%v", value)
	}

	err = flag.Prologue()
	return
}

func (flag *Flag) Prologue() (err error) {
	switch flag.Value.Interface().(type) {
	case time.Time, net.IP, net.IPNet, net.Interface:
	case *time.Time, *os.File, *net.IP, *net.IPNet, *net.Interface:
	default:
		if err = flag.prologue(flag.Value.Type()); err != nil {
			err = fmt.Errorf("cannot set %v as flag: %v", flag.StructField.Name, err)
			return
		}
	}

	return
}

func (flag *Flag) prologue(typ reflect.Type) (err error) {
	switch typ.Kind() {
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
		err = flag.prologue(typ.Elem())
		return
	default:
		return
	}

	err = fmt.Errorf("not support type: %v", typ.Kind())
	return
}

// parse the pass argument, should consumed one and only one argument
func (flag *Flag) Parse(args ...string) (n int, err error) {
	if len(args) == 0 {
		err = fmt.Errorf("should pass %v", flag.Hint())
		return
	}

	if len(flag.choise) > 0 {
		found := false

		for _, token := range flag.choise {
			if token == args[0] {
				found = true
				break
			}
		}

		if !found {
			err = fmt.Errorf("should pass one of [%v]: %v", strings.Join(flag.choise, " "), args[0])
			return
		}
	}

	// the special case
	switch flag.Value.Interface().(type) {
	case time.Duration, *time.Duration:
		var duration time.Duration

		if duration, err = time.ParseDuration(args[0]); err == nil {
			flag.setValue(reflect.ValueOf(&duration))
			n++
			return
		}
	case time.Time, *time.Time:
		var timestamp time.Time

		if timestamp, err = time.Parse(time.RFC3339, args[0]); err == nil {
			flag.setValue(reflect.ValueOf(&timestamp))
			n++
			return
		}

		err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
		return
	case os.File, *os.File:
		var file *os.File

		if file, err = os.Open(args[0]); err == nil {
			flag.setValue(reflect.ValueOf(file))
			n++
			return
		}

		if file != nil {
			flag.setValue(reflect.ValueOf(file))
			n++
			return
		}

		err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
		return
	case net.IP, *net.IP:
		var ip net.IP

		if ip = net.ParseIP(args[0]); ip == nil {
			ips, err := net.LookupIP(args[0])
			if err == nil && len(ips) > 0 {
				ip = ips[0]
			}
		}

		if ip != nil {
			flag.setValue(reflect.ValueOf(&ip))
			n++
			return
		}

		err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
		return
	case net.IPNet, *net.IPNet:
		var inet *net.IPNet

		_, inet, err = net.ParseCIDR(args[0])

		if inet != nil {
			flag.setValue(reflect.ValueOf(inet))
			n++
			return
		}

		err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
		return
	case net.Interface, *net.Interface:
		var iface *net.Interface

		if iface, err = net.InterfaceByName(args[0]); err == nil {
			flag.setValue(reflect.ValueOf(iface))
			n++
			return
		}

		err = fmt.Errorf("should pass %v: %v", flag.Hint(), args[0])
		return
	}

	if err = flag.parse(flag.Value, args[0]); err != nil {
		// cannot parse the built-in type
		return
	}

	n++
	return
}

func (flag *Flag) parse(value reflect.Value, args string) (err error) {
	switch kind := value.Type().Kind(); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var v int

		if v, err = strconv.Atoi(args); err != nil {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args)
			return
		}

		value.SetInt(int64(v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var v int

		if v, err = strconv.Atoi(args); err != nil {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args)
			return
		} else if v < 0 {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args)
			return
		}

		value.SetUint(uint64(v))
	case reflect.Float32, reflect.Float64:
		rat := &big.Rat{}
		if _, ok := rat.SetString(args); !ok {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args)
			return
		}

		float, exact := rat.Float64()
		flag.Infof("convert %v to %v (exact: %v)", args, float, exact)
		value.SetFloat(float)
	case reflect.Complex64, reflect.Complex128:
		var cplx complex128

		if cplx, err = strconv.ParseComplex(args, 128); err != nil {
			err = fmt.Errorf("should pass %v: %v", flag.Hint(), args)
			return
		}
		value.SetComplex(cplx)
	case reflect.String:
		// copy the argument to string
		value.SetString(args)
	case reflect.Ptr:
		shadow := reflect.New(value.Type().Elem())
		if err = flag.parse(shadow.Elem(), args); err != nil {
			// cannot setup the value
			return
		}
		value.Set(shadow)
	case reflect.Slice:
		shadow := reflect.New(value.Type().Elem()).Elem()
		if err = flag.parse(shadow, args); err != nil {
			// cannot setup the value
			return
		}
		value.Set(reflect.Append(value, shadow))
	default:
		err = fmt.Errorf("not support parse %v: %v", kind, value)
		return
	}

	return
}

func (flag *Flag) setValue(value reflect.Value) {
	switch flag.Value.Kind() {
	case reflect.Ptr:
		flag.Value.Set(value)
	default:
		flag.Value.Set(value.Elem())
	}
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
	hint = flag.hint(flag.StructField.Type)
	return
}

func (flag *Flag) hint(typ reflect.Type) (hint string) {
	switch kind := typ.Kind(); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		hint = "INT"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		hint = "UINT"
	case reflect.Float32, reflect.Float64:
		hint = "RAT"
	case reflect.Complex64, reflect.Complex128:
		hint = "CPLX"
	case reflect.String:
		hint = "STR"
	case reflect.Ptr:
		hint = flag.hint(typ.Elem())
	case reflect.Slice:
		hint = fmt.Sprintf("[%v ...]", flag.hint(typ.Elem()))
	case reflect.Array:
		hint = fmt.Sprintf("[%v %v]", flag.hint(typ.Elem()), typ.Len())
	default:
		switch flag.Value.Interface().(type) {
		case time.Duration, *time.Duration:
			hint = "TIME"
		case os.File, *os.File:
			hint = "FILE"
		case net.IP, *net.IP:
			hint = "IP"
		case net.IPNet, *net.IPNet:
			hint = "CIDR"
		case net.Interface, *net.Interface:
			hint = "IFACE"
		}
	}

	return
}

// set the choise value
func (flag *Flag) SetChoice(choise []string) (err error) {
	flag.choise = append(flag.choise, choise...)
	return
}

// get the choice
func (flag *Flag) GetChoice() (choise []string) {
	choise = append(choise, flag.choise...)
	return
}

// the default value
func (flag *Flag) Default() (_default string) {
	_default = flag._default
	return
}
