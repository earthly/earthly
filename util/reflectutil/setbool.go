package reflectutil

import (
	"reflect"
)

// SetBool looks for a boolean field value named `fieldName` in `iface` and sets it to `value`.
// Upon success a true value is returned, otherwise false.
func SetBool(iface interface{}, fieldName string, value bool) bool {
	rv := reflect.ValueOf(iface)
	for rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	field := rv.FieldByName(fieldName)
	if field.IsValid() && field.Bool() {
		field.SetBool(value)
		return true
	}
	return false
}
