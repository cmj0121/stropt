package stropt

import (
	"fmt"
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

// the callback function called when field set.
type Callback func(stropt *StrOpt, field Field)

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

	callback, ok := callbacks_pool[name]
	if !ok {
		err = fmt.Errorf("callback not found: %v", name)
		return
	}

	callback(stropt, field)
	return
}
