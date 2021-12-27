package stropt

import (
	"reflect"
)

// the settable field in the stropt which used to show the option meta
// and validate the pass arguments
type Field interface {
	// parse and fill the field, return number of args used or error when failure
	Parse(args ...string) (n int, err error)

	// name of the field
	GetName() string
	// shortcut of the field
	GetShortcut() string
	// the customized tag of the field
	GetTag() reflect.StructTag
}
