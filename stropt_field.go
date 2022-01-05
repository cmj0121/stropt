package stropt

import (
	"fmt"
	"reflect"
)

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

// return the shortcut of the field
func (stropt *StrOpt) GetShortcut() (shortcut string) {
	// sub-command does not contains shortcut
	return
}

// set the choise value
func (stropt *StrOpt) SetChoice(choise []string) (err error) {
	err = fmt.Errorf("sub-command should not set choise")
	return
}

// get the choice
func (stropt *StrOpt) GetChoice() (choise []string) {
	return
}

// the default value
func (stropt *StrOpt) Default() (def string) {
	return
}
