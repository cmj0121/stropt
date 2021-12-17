package stropt

import (
	"fmt"
	"reflect"
)

// the instance of stropt which serves the *Struct, ready to parse the
// input arguments and fill the data into Struct.
type StrOpt struct {
	reflect.Value
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
