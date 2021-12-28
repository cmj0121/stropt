package stropt

import (
	"errors"
	"fmt"
	"reflect"
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
	CALLBACK_HELP = "Help_"
	// the pre-defined callback, show the version info
	CALLBACK_VERSION = "Version_"
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

	if stropt == nil {
		err = fmt.Errorf("should provides valid stropt: %v", stropt)
		return
	}

	stropt.Debugf("try call callback %#v", name)

	// call local callback if exists
	callback_value := stropt.Value.MethodByName(name)
	if callback_value.IsValid() && !callback_value.IsZero() {
		// since the method is not Callback, need to convert type to function ptr
		callback, ok := callback_value.Interface().(func(stropt *StrOpt, field Field) error)
		if ok {
			stropt.Infof("call local callback: %v", name)
			// found the callback, call it
			if err = callback(stropt, field); err != ERR_CALLBACK_NOT_IMPLEMENTED {
				// only return when callback implemented
				return
			}
		}
	}

	// call global callback if exists
	value := reflect.ValueOf(stropt)
	callback_value = value.MethodByName(name)

	if callback_value.IsValid() && !callback_value.IsZero() {
		// since the method is not Callback, need to convert type to function ptr
		callback, ok := callback_value.Interface().(func(stropt *StrOpt, field Field) error)
		if ok {
			stropt.Infof("call local callback: %v", name)
			// found the callback, call it
			if err = callback(stropt, field); err != ERR_CALLBACK_NOT_IMPLEMENTED {
				// only return when callback implemented
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

	stropt.Infof("call global callback: %v", name)
	err = callback(stropt, field)
	return
}
