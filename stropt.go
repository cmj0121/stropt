package stropt

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/cmj0121/trace"
)

// the instance of stropt which serves the *Struct, ready to parse the
// input arguments and fill the data into Struct.
type StrOpt struct {
	reflect.Value

	// the shadow of the value, create and copy to original value
	shadow reflect.Value

	// the log sub-system
	*trace.Tracer

	// name of the stropt, always be the lowercase
	name string

	// the tag of the StrOpt
	tag reflect.StructTag
	// the flip/flag fields, append by the order
	fields []Field
	// named field, should be flip/flag
	named_fields map[string]Field
	// the arguments field, trigger sequencity
	args_fields []Field
	args_idx    int
	// the named sub-command
	sub_fields map[string]Field
	// the version info
	version string
}

// create an instance of StrOpt by input *StrOpt, may return error
// if input value is invalid.
func New(in interface{}) (stropt *StrOpt, err error) {
	value := reflect.ValueOf(in)
	value.Kind()

	// only allow pass the *Struct
	if !(value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
		err = fmt.Errorf("should pass *Struct: %T", in)
		return
	}

	stropt = &StrOpt{
		Value: value,
		// setup the tracer
		Tracer: trace.GetTracer(PROJ_NAME),
		// the internal fields
		name: strings.ToLower(value.Elem().Type().Name()),
		// named fields, call when parse argument and set field
		named_fields: map[string]Field{},
		// sub-command fields
		sub_fields: map[string]Field{},
	}

	stropt.Infof("new StrOpt: %[1]v (%[1]T)", in)
	// pass the type of Struct (not the *Struct)
	err = stropt.prologue(stropt.Value.Elem(), reflect.TypeOf(in).Elem())
	return
}

// the helper function for create StrOpt, raise panic if catch error.
func MustNew(in interface{}) (stropt *StrOpt) {
	var err error

	if stropt, err = New(in); err != nil {
		panic(err)
	}
	return
}

// override the name of the command-line usage, and always be the lowercase.
func (stropt *StrOpt) Name(name string) {
	// set the name as lower-case
	stropt.Tracef("change StrOpt name: %#v", name)
	stropt.name = strings.ToLower(name)
}

// write the usage to the pass io.Writer
func (stropt *StrOpt) Usage(w io.Writer) {
	buff := &bytes.Buffer{}

	var usage []string

	switch {
	case len(stropt.fields) > 0 && len(stropt.args_fields) > 0:
		usage = append(usage, fmt.Sprintf("usage: %v [OPTION] [ARGS] ...", stropt.name))
		usage = append(usage, "")
	case len(stropt.fields) > 0:
		usage = append(usage, fmt.Sprintf("usage: %v [OPTION]", stropt.name))
		usage = append(usage, "")
	case len(stropt.args_fields) > 0:
		usage = append(usage, fmt.Sprintf("usage: %v [ARGS] ...", stropt.name))
		usage = append(usage, "")
	default:
		usage = append(usage, fmt.Sprintf("usage: %v", stropt.name))
		usage = append(usage, "")
	}

	if len(stropt.fields) > 0 {
		usage = append(usage, "options:")
		for _, field := range stropt.fields {
			usage = append(usage, stropt.description(field, false))
		}
		usage = append(usage, "")
	}

	if len(stropt.args_fields) > 0 {
		usage = append(usage, "arguments:")
		for _, field := range stropt.args_fields {
			usage = append(usage, stropt.description(field, true))
		}
		usage = append(usage, "")
	}

	if len(stropt.sub_fields) > 0 {
		usage = append(usage, "sub-commands:")
		for _, field := range stropt.sub_fields {
			usage = append(usage, stropt.description(field, true))
		}
		usage = append(usage, "")
	}

	buff.WriteString(strings.Join(usage, "\n"))
	w.Write(buff.Bytes()) // nolint
}

// parse the input arguments and fill the *Struct, return error when failure.
func (stropt *StrOpt) Parse(args ...string) (n int, err error) {
	stropt.Tracef("start parse: %v", args)
	defer func() {
		if stropt.shadow.IsValid() && !stropt.shadow.IsZero() {
			// copy the shadow value to current value
			stropt.Value.Set(stropt.shadow)
		}
	}()

	no_option := false
	idx := 0
	for idx < len(args) {
		nargs := 0
		token := args[idx]
		stropt.Debugf("parse #%v: %v", idx, token)

		switch {
		case token == "--":
			no_option = true
			stropt.Infof("explicit claims no options remains")
		case !no_option && len(token) > 2 && token[:2] == "--":
			field, ok := stropt.named_fields[token[2:]]
			if !ok {
				err = fmt.Errorf("option %v not found", token)
				return
			} else if nargs, err = stropt.parse(field, args[idx+1:]...); err != nil {
				err = fmt.Errorf("parse %v fail: %v", token, err)
				return
			}

			idx += nargs
		case !no_option && len(token) > 1 && token[:1] == "-":
			switch len(token) {
			case 2:
				// single shortcut
				field, ok := stropt.named_fields[token[1:]]
				if !ok {
					err = fmt.Errorf("option %v not found", token)
					return
				} else if nargs, err = stropt.parse(field, args[idx+1:]...); err != nil {
					err = fmt.Errorf("parse %v fail: %v", token, err)
					return
				}

				idx += nargs
			default:
				// multiple shortcut
				for _, shortcut := range token[1:] {
					field, ok := stropt.named_fields[string(shortcut)]
					if !ok {
						err = fmt.Errorf("option -%v not found", token)
						return
					} else if _, err = stropt.parse(field); err != nil {
						err = fmt.Errorf("parse -%v fail: %v", shortcut, err)
						return
					}
				}
			}
		default:
			switch field, ok := stropt.sub_fields[token]; ok {
			case true:
				// sub-command
				if _, err = stropt.parse(field, args[idx+1:]...); err != nil {
					err = fmt.Errorf("parse %v fail: %v", token, err)
					return
				}

				return
			case false:
				// position field
				if stropt.args_idx >= len(stropt.args_fields) {
					err = fmt.Errorf("unknown argument: %v", token)
					return
				}

				field = stropt.args_fields[stropt.args_idx]
				if nargs, err = stropt.parse(field, args[idx:]...); err != nil {
					err = fmt.Errorf("parse %v fail: %v", token, err)
					return
				}

				stropt.args_idx++
			}

			// note the args already take-out the first argument, which
			// counts when break the switch-statement
			idx += nargs - 1
		}

		idx++
	}

	return
}

// the helper utility for parse the arguments and trigger callback with specified field
func (stropt *StrOpt) parse(field Field, args ...string) (n int, err error) {
	stropt.Debugf("parse %v on %v", args, field)
	if n, err = field.Parse(args...); err == nil {
		if name, ok := field.GetTag().Lookup(KEY_CALLBACK); ok {
			// call the callback function
			err = CallCallback(name, stropt, field)
		}
	}

	return
}

// parse from the command-line arguments
func (stropt *StrOpt) Run() {
	if _, err := stropt.Parse(os.Args[1:]...); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("error: %v\n", err))
		stropt.Usage(os.Stderr)
		os.Exit(1)
	}
}

// parse the pass Struct into fields
func (stropt *StrOpt) prologue(value reflect.Value, typ reflect.Type) (err error) {
	stropt.Debugf("stropt prologue %v (%v)", value, typ)

	for idx := 0; idx < typ.NumField(); idx++ {
		field_type := typ.Field(idx)
		field_value := value.Field(idx)

		var field Field
		stropt.Tracef("process #%d field: %v (%v, %v)", idx, field_value, field_type, field_type.Tag)
		if field, err = stropt.setField(field_value, field_type); err != nil {
			err = fmt.Errorf("set #%v field: %v", idx, err)
			return
		}

		if value, ok := field_type.Tag.Lookup(KEY_CHOICE); ok {
			choice := strings.Split(value, " \t")
			if err = field.SetChoice(choice); err != nil {
				err = fmt.Errorf("set #%v field: %v", idx, err)
				return
			}
		}
	}
	return
}

// set the pass reflect.Value and reflect.StructField to Field
func (stropt *StrOpt) setField(value reflect.Value, typ reflect.StructField) (field Field, err error) {
	force_as_flag := false
	if v, ok := typ.Tag.Lookup(KEY_ATTR); ok {
		attrs := strings.Split(v, " \t")
		sort.Strings(attrs)

		idx := sort.SearchStrings(attrs, KEY_ATTR_FLAG)
		if idx >= 0 && idx < len(attrs) && attrs[idx] == KEY_ATTR_FLAG {
			stropt.Infof("attribute: %v", KEY_ATTR_FLAG)
			force_as_flag = true
		}
	}

	switch {
	case !value.CanSet():
		stropt.Debugf("field #%v cannot be set, skip", value)
		return
	case strings.TrimSpace(string(typ.Tag)) == TAG_IGNORE:
		stropt.Debugf("field %v expressily been skip", value)
		return
	default:
		stropt.Debugf("set field %v (%v)", value, typ.Type.Kind())

		// specified case
		switch value.Interface().(type) {
		case time.Time, net.IP, net.IPNet, net.Interface:
			if field, err = NewFlag(stropt.Tracer, value, typ); err != nil {
				err = fmt.Errorf("new flag from %v: %v", value, err)
				return
			}
			err = stropt.setOption(field)
			return
		case *time.Time, *os.File, *net.IP, *net.IPNet, *net.Interface:
			switch force_as_flag {
			case true:
				if field, err = NewFlag(stropt.Tracer, value, typ); err != nil {
					err = fmt.Errorf("new flag from %v: %v", value, err)
					return
				}
				err = stropt.setOption(field)
			case false:
				if field, err = NewArgument(stropt.Tracer, value, typ); err != nil {
					err = fmt.Errorf("new flag from %v: %v", value, err)
					return
				}
				err = stropt.setArgument(field)
			}
			return
		}

		switch typ.Type.Kind() {
		case reflect.Ptr: // argument or sub-command
			raw_type := typ.Type.Elem()

			switch {
			case raw_type.Kind() == reflect.Struct: // sub-command
				// create a new StrOpt as the sub-command
				name := strings.ToLower(typ.Name)
				if value, ok := typ.Tag.Lookup(KEY_NAME); ok {
					stropt.Debugf("override the field name %#v: %#v", name, value)
					name = strings.ToLower(value)
				}

				shadow := reflect.New(typ.Type.Elem())

				sub := &StrOpt{
					Value:        value,
					shadow:       shadow,
					Tracer:       stropt.Tracer,
					name:         name,
					tag:          typ.Tag,
					named_fields: map[string]Field{},
				}

				if err = sub.prologue(shadow.Elem(), typ.Type.Elem()); err != nil {
					err = fmt.Errorf("cannot set sub-command: %v", err)
					return
				}

				err = stropt.setSub(sub)
			case force_as_flag:
				if field, err = NewFlag(stropt.Tracer, value, typ); err != nil {
					err = fmt.Errorf("new flag from %v: %v", value, err)
					return
				}
				err = stropt.setOption(field)
			default:
				if field, err = NewArgument(stropt.Tracer, value, typ); err != nil {
					err = fmt.Errorf("new flag from %v: %v", value, err)
					return
				}
				err = stropt.setArgument(field)
			}
		case reflect.Struct: // may embedded field
			err = stropt.prologue(value, typ.Type)
			return
		case reflect.Bool: // the flip option
			if field, err = NewFlip(stropt.Tracer, value, typ); err != nil {
				err = fmt.Errorf("new flip from %v: %v", value, err)
				return
			}
			err = stropt.setOption(field)
		default: // may flag option
			if field, err = NewFlag(stropt.Tracer, value, typ); err != nil {
				err = fmt.Errorf("new flag from %v: %v", value, err)
				return
			}
			err = stropt.setOption(field)
		}
	}

	return
}

// set as the named (name/shortcut) options
func (stropt *StrOpt) setOption(field Field) (err error) {
	stropt.fields = append(stropt.fields, field)

	// set named field
	name := field.GetName()
	if _, ok := stropt.named_fields[name]; ok {
		err = fmt.Errorf("duplicate field name: %v", name)
		return
	}
	stropt.named_fields[name] = field

	if shortcut := field.GetShortcut(); len(shortcut) > 0 {
		if _, ok := stropt.named_fields[shortcut]; ok {
			err = fmt.Errorf("duplicate field shortcut: %v", shortcut)
			return
		}
		stropt.named_fields[shortcut] = field
	}

	if _default, ok := field.GetTag().Lookup(KEY_DEFAULT); ok {
		// set the default value before parse
		_, err = field.Parse(_default)
	}

	return
}

// set as the arguments
func (stropt *StrOpt) setArgument(field Field) (err error) {
	stropt.args_fields = append(stropt.args_fields, field)

	if _default, ok := field.GetTag().Lookup(KEY_DEFAULT); ok {
		// set the default value before parse
		_, err = field.Parse(_default)
	}
	return
}

// set as the sub-command
func (stropt *StrOpt) setSub(field Field) (err error) {
	name := field.GetName()
	if _, ok := stropt.sub_fields[name]; ok {
		err = fmt.Errorf("duplicate sub-command: %v", name)
		return
	}
	stropt.sub_fields[name] = field
	return
}

// show the field description
func (stropt *StrOpt) description(field Field, sub bool) (desc string) {
	name := field.GetName()
	shortcut := field.GetShortcut()

	switch {
	case sub:
	case len(name) > 0 && len(shortcut) > 0:
		shortcut = fmt.Sprintf("-%v", shortcut)
		name = fmt.Sprintf("--%v", name)
	case len(name) > 0:
		name = fmt.Sprintf("--%v", name)
	case len(shortcut) > 0:
		shortcut = fmt.Sprintf("-%v", shortcut)
	}

	// the helper description of the field
	help, _ := field.GetTag().Lookup(KEY_DESC)

	choice := field.GetChoice()
	if len(choice) > 0 {
		// append the choice
		help = fmt.Sprintf("%v [%v]", help, strings.Join(choice, " "))
	}

	desc = fmt.Sprintf("%3v %v", shortcut, name)
	switch _default := field.Default(); _default {
	case "":
		desc = fmt.Sprintf("    %-22v %v", desc, help)
	default:
		desc = fmt.Sprintf("    %-22v %v [default: %v]", desc, help, _default)
	}

	desc = strings.TrimRight(desc, " ")
	return
}

// show the version info
func (stropt *StrOpt) Version_(_stropt *StrOpt, _field Field) (err error) {
	switch stropt.version {
	case "":
		err = ERR_CALLBACK_NOT_IMPLEMENTED
	default:
		os.Stdout.WriteString(stropt.version)
		os.Exit(0)
	}
	return
}

func (stropt *StrOpt) Version(ver string) {
	stropt.version = ver
}
