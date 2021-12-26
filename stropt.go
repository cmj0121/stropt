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

	// the flip/flag fields, append by the order
	fields []Field
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
	err = stropt.prologue(reflect.TypeOf(in).Elem())
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

	usage := []string{}

	switch {
	case len(stropt.fields) > 0:
		usage = append(usage, fmt.Sprintf("usage: %v [OPTION]", stropt.name))
		usage = append(usage, "")

		usage = append(usage, "options:")
		for _, field := range stropt.fields {
			usage = append(usage, field.Description())
		}
	default:
		usage = append(usage, fmt.Sprintf("usage: %v", stropt.name))
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
func (stropt *StrOpt) prologue(typ reflect.Type) (err error) {
	for idx := 0; idx < typ.NumField(); idx++ {
		field_type := typ.Field(idx)
		field_value := stropt.Value.Elem().Field(idx)

		var field Field
		stropt.Tracef("process #%d field: %v (%v)", idx, field_value, field_type)

		switch {
		case !field_value.CanSet():
			stropt.Debugf("field #%v cannot be set, skip", idx)
		case strings.TrimSpace(string(field_type.Tag)) == TAG_IGNORE:
			stropt.Debugf("field #%v expressily been skip", idx)
		default:
			if field, err = stropt.setField(field_value, field_type); err != nil {
				err = fmt.Errorf("set #%v field: %v", idx, err)
				return
			}
		}

		// modify the field attributes
		stropt.modifyField(field, field_type)
		stropt.fields = append(stropt.fields, field)
	}
	return
}

// set the pass reflect.Value and reflect.StructField to Field
func (stropt *StrOpt) setField(value reflect.Value, typ reflect.StructField) (field Field, err error) {
	switch typ.Type.Kind() {
	case reflect.Ptr: // argument or sub-command
		err = fmt.Errorf("not implement %v (%v)", value, typ)
		return
	case reflect.Struct: // may embedded field
		err = fmt.Errorf("not implement %v (%v)", value, typ)
		return
	case reflect.Bool: // the flip option
		field, err = NewFlip(stropt.Tracer, value)
	default: // may flag option
		err = fmt.Errorf("not implement %v (%v)", value, typ)
		return
	}
	return
}

// setup / modify the field attributes
func (stropt *StrOpt) modifyField(field Field, typ reflect.StructField) {
	desc := fmt.Sprintf("    --%v", strings.ToLower(typ.Name))
	field.SetDescription(desc)
}
