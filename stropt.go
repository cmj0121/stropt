package stropt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// the instance of stropt which serves the *Struct, ready to parse the
// input arguments and fill the data into Struct.
type StrOpt struct {
	reflect.Value

	// name of the stropt, always be the lowercase
	name string
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
		// the internal fields
		name: strings.ToLower(value.Elem().Type().Name()),
	}

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
	stropt.name = strings.ToLower(name)
}

// write the usage to the pass io.Writer
func (stropt *StrOpt) Usage(w io.Writer) {
	buff := &bytes.Buffer{}

	buff.WriteString(fmt.Sprintf("usage: %v\n", stropt.name))

	w.Write(buff.Bytes()) // nolint
}

// parse the input arguments and fill the *Struct, return error when failure.
func (stropt *StrOpt) Parse(args ...string) (err error) {
	stropt.Usage(os.Stderr)
	return
}
