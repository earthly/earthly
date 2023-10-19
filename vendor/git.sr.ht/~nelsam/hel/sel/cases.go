package sel

import (
	"fmt"
	"reflect"
)

// Cases takes either a mock struct (i.e. a struct of only channel fields)
// or a channel and returns a []reflect.SelectCase that may be used in
// reflect.Select calls. If dir is reflect.SelectSend, then send must contain
// exactly one value for each channel in mock.
func Cases(dir reflect.SelectDir, mock interface{}, send ...interface{}) ([]reflect.SelectCase, error) {
	v := reflect.ValueOf(mock)
	switch v.Kind() {
	case reflect.Chan:
		sel := reflect.SelectCase{Dir: dir, Chan: v}
		if dir != reflect.SelectSend {
			return []reflect.SelectCase{sel}, nil
		}
		if len(send) != 1 {
			return nil, fmt.Errorf("sel: expected 1 argument for mock (%v); got %d", v.Type(), len(send))
		}
		arg := send[0]
		argV := reflect.ValueOf(arg)
		if arg == nil {
			argV = reflect.Zero(v.Type().Elem())
		}
		if !argV.Type().ConvertibleTo(v.Type().Elem()) {
			return nil, fmt.Errorf("sel: argument type %v is not convertible to mock type %v", argV.Type(), v.Type().Elem())
		}
		sel.Send = argV.Convert(v.Type().Elem())
		return []reflect.SelectCase{sel}, nil
	case reflect.Struct:
		if v.NumField() == 0 {
			return nil, fmt.Errorf("sel: empty struct (%v) is not a valid mock field", v.Type())
		}
		switch dir {
		case reflect.SelectSend:
			if len(send) != v.NumField() {
				argString := "argument"
				if v.NumField() != 1 {
					argString = "arguments"
				}
				return nil, fmt.Errorf("sel: expected %d %s for mock (%v); got %d", v.NumField(), argString, v.Type(), len(send))
			}
			var cases []reflect.SelectCase
			for i := 0; i < v.NumField(); i++ {
				c, err := Cases(dir, v.Field(i).Interface(), send[i])
				if err != nil {
					return nil, fmt.Errorf("sel: field %d of type %v cannot be used as a mock field: %w", i, v.Type(), err)
				}
				cases = append(cases, c...)
			}
			return cases, nil
		default:
			var cases []reflect.SelectCase
			for i := 0; i < v.NumField(); i++ {
				c, err := Cases(dir, v.Field(i).Interface())
				if err != nil {
					return nil, fmt.Errorf("sel: field %d of type %v cannot be used as a mock field: %w", i, v.Type(), err)
				}
				cases = append(cases, c...)
			}
			return cases, nil
		}
	default:
		return nil, fmt.Errorf("sel: type %v is neither a struct of channels nor a channel", v.Type())
	}
}
