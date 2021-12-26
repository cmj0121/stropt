package stropt

// the settable field in the stropt which used to show the option meta
// and validate the pass arguments
type Field interface {
	// parse and fill the field, return number of args used or error when failure
	Parse(args ...string) (n int, err error)

	// set the field description
	SetDescription(string)
	// get the field description
	Description() string
}
