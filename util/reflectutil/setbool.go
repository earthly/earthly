package reflectutil

import (
	"reflect"
)

// SetBool looks for a boolean struct value named `structName` in `iface` and sets it to `value`.
// Upon success a true value is returned, otherwise false.
func SetBool(iface interface{}, structName string, value bool) bool {
	rv := reflect.ValueOf(iface)
	for rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	field := rv.FieldByName(structName)
	if field.IsValid() && field.Bool() {
		field.SetBool(value)
		return true
	}
	return false
}
