package stropt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/cmj0121/trace"
)

// the instance of stropt which serves the *Struct, ready to parse the
// input arguments and fill the data into Struct.
type StrOpt struct {
	reflect.Value

	// the log sub-system
	*trace.Tracer

	// name of the stropt, always be the lowercase
	name string

	// the tag of the StrOpt
	tag reflect.StructTag
	// the flip/flag fields, append by the order
	fields []Field
	// the sub-commands
	subs []Field
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
	case len(stropt.fields) > 0 && len(stropt.subs) > 0:
		usage = append(usage, fmt.Sprintf("usage: %v [OPTION] [SUB-COMMAND]", stropt.name))
		usage = append(usage, "")
	case len(stropt.fields) > 0:
		usage = append(usage, fmt.Sprintf("usage: %v [SUB-COMMAND]", stropt.name))
		usage = append(usage, "")
	case len(stropt.subs) > 0:
		usage = append(usage, fmt.Sprintf("usage: %v [OPTION]", stropt.name))
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

	if len(stropt.subs) > 0 {
		usage = append(usage, "sub-commands:")
		for _, field := range stropt.subs {
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
	stropt.Usage(os.Stderr)
	return
}

// parse the pass Struct into fields
func (stropt *StrOpt) prologue(value reflect.Value, typ reflect.Type) (err error) {
	stropt.Debugf("stropt prologue %v (%v)", value, typ)

	for idx := 0; idx < typ.NumField(); idx++ {
		field_type := typ.Field(idx)
		field_value := value.Field(idx)

		stropt.Tracef("process #%d field: %v (%v, %v)", idx, field_value, field_type, field_type.Tag)
		if err = stropt.setField(field_value, field_type); err != nil {
			err = fmt.Errorf("set #%v field: %v", idx, err)
			return
		}
	}
	return
}

// set the pass reflect.Value and reflect.StructField to Field
func (stropt *StrOpt) setField(value reflect.Value, typ reflect.StructField) (err error) {
	var field Field

	switch {
	case !value.CanSet():
		stropt.Debugf("field #%v cannot be set, skip", value)
		return
	case strings.TrimSpace(string(typ.Tag)) == TAG_IGNORE:
		stropt.Debugf("field %v expressily been skip", value)
		return
	default:
		stropt.Debugf("set field %v (%v)", value, typ.Type.Kind())

		switch typ.Type.Kind() {
		case reflect.Ptr: // argument or sub-command
			// create a new StrOpt as the sub-command
			name := strings.ToLower(typ.Name)
			if value, ok := typ.Tag.Lookup(KEY_NAME); ok {
				stropt.Debugf("override the field name %#v: %#v", name, value)
				name = strings.ToLower(value)
			}

			sub := &StrOpt{
				Value:  value,
				Tracer: stropt.Tracer,
				name:   name,
				tag:    typ.Tag,
				fields: []Field{},
			}

			new_value := reflect.New(typ.Type.Elem()).Elem()
			if err = sub.prologue(new_value, typ.Type.Elem()); err != nil {
				err = fmt.Errorf("cannot set sub-command: %v", err)
				return
			}

			stropt.subs = append(stropt.subs, sub)
			return
		case reflect.Struct: // may embedded field
			err = stropt.prologue(value, typ.Type)
			return
		case reflect.Bool: // the flip option
			if field, err = NewFlip(stropt.Tracer, value, typ); err != nil {
				err = fmt.Errorf("new flip from %v: %v", value, err)
				return
			}
		default: // may flag option
			if field, err = NewFlag(stropt.Tracer, value, typ); err != nil {
				err = fmt.Errorf("new flag from %v: %v", value, err)
				return
			}
		}
	}

	stropt.fields = append(stropt.fields, field)
	return
}

// show the field description
func (stropt *StrOpt) description(field Field, sub bool) (desc string) {
	var shortcut string
	name := strings.ToLower(field.GetName())

	if value, ok := field.GetTag().Lookup(KEY_NAME); ok {
		// override the field's name
		stropt.Tracef("override field name %v: %v", name, value)
		name = strings.ToLower(value)
	}

	if value, ok := field.GetTag().Lookup(KEY_SHORTCUT); ok && !sub {
		// override the field' shortcut, should be rune
		runes := []rune(value)
		switch {
		case len(runes) == 0:
			// no changed
		case len(runes) > 1:
			stropt.Warnf("shortcut too large: %v (should be one and only one rune", value)
		default:
			shortcut = string(runes[0])
		}
	}

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

	desc = fmt.Sprintf("%3v %v", shortcut, name)
	desc = strings.TrimRight(fmt.Sprintf("    %-22v %v", desc, help), " ")
	return
}

// return the Tag of the field
func (stropt *StrOpt) GetTag() (tag reflect.StructTag) {
	tag = stropt.tag
	return
}

// return the original name of the field
func (stropt *StrOpt) GetName() (name string) {
	name = stropt.name
	return
}
