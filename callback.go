package stropt

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	// the global callbacks, register explicit
	callbacks_pool = map[string]Callback{}
	// the global lock when register callback
	callbacks_lock = sync.Mutex{}
)

// pre-defined callback function
var (
	// the pre-defined callback, show the help message
	CALLBACK_HELP = "_help"
	// the pre-defined callback, show the version info
	CALLBACK_VERSION = "_version"
)

var (
	// pre-defined error for not implemenetd callback
	ERR_CALLBACK_NOT_IMPLEMENTED = errors.New("callback not implemenetd")
)

// the callback function called when field set.
type Callback func(stropt *StrOpt, field Field) error

func RegisterCallback(name string, callback Callback) {
	callbacks_lock.Lock()
	defer callbacks_lock.Unlock()

	if _, ok := callbacks_pool[name]; ok {
		// duplicated callback, raise panic
		panic(fmt.Sprintf("duplicate callback: %v", name))
	}

	callbacks_pool[name] = callback
}

func CallCallback(name string, stropt *StrOpt, field Field) (err error) {
	callbacks_lock.Lock()
	defer callbacks_lock.Unlock()

	if stropt != nil {
		// call local callback if exists
		value := reflect.ValueOf(stropt).Elem()
		// always convert the name to their Unicode title case
		callback_value := value.MethodByName(strings.ToTitle(name))

		if !callback_value.IsZero() {
			callback, ok := callback_value.Interface().(Callback)
			if ok {
				// found the callback, call it
				err = callback(stropt, field)
				return
			}
		}
	}

	// call global callback if exists
	callback, ok := callbacks_pool[name]
	if !ok {
		err = fmt.Errorf("callback not found: %v", name)
		return
	}

	err = callback(stropt, field)
	return
}
