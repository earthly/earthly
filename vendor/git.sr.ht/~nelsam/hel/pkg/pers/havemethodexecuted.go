// This is free and unencumbered software released into the public
// domain.  For more information, see <http://unlicense.org> or the
// accompanying UNLICENSE file.

package pers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/poy/onpar/diff"
	"github.com/poy/onpar/expect"
	"github.com/poy/onpar/matchers"
)

// Matcher is any type that can match values.  Some code in this package supports
// matching against child matchers, for example:
//
//	HaveBeenExecuted("Foo", WithArgs(matchers.HaveLen(12)))
type Matcher interface {
	Match(actual interface{}) (interface{}, error)
}

type any int

// Any is a special value to tell pers to allow any value at the position used.
// For example, you can assert only on the second argument with:
//
//	HaveMethodExecuted("Foo", WithArgs(Any, 22))
const Any any = -1

type variadicAny int

// VariadicAny is a special value, similar to Any, but specifically to tell pers
// to allow any number of values for the variadic arguments. It must be passed
// in after all non-variadic arguments so that its position matches the position
// that variadic arguments would normally be.
//
// This cannot be used to check some variadic arguments without checking the
// others - at least right now, you must either assert on all of the variadic
// arguments or none of them.
const VariadicAny variadicAny = -2

// HaveMethodExecutedOption is an option function for the HaveMethodExecutedMatcher.
type HaveMethodExecutedOption func(HaveMethodExecutedMatcher) HaveMethodExecutedMatcher

// Within returns a HaveMethodExecutedOption which sets the HaveMethodExecutedMatcher
// to be executed within a given timeframe.
func Within(d time.Duration) HaveMethodExecutedOption {
	return func(m HaveMethodExecutedMatcher) HaveMethodExecutedMatcher {
		m.within = d
		return m
	}
}

// WithArgs returns a HaveMethodExecutedOption which sets the HaveMethodExecutedMatcher
// to only pass if the latest execution of the method called it with the passed in
// arguments.
//
// WithArgs will cause a panic if:
//   - The argument passed in is not ConvertibleTo the mock's argument type at the
//     matching index.
//   - The number of arguments does not match the number of arguments in the mock.
//   - For variadic methods, the number of arguments in the mock is a minimum.
//     If len(args) is >= len(nonVariadicArgs(mock)) but len(args) !=
//     len(totalArgs(mock)), the matcher will return an error instead.
func WithArgs(args ...interface{}) HaveMethodExecutedOption {
	return func(m HaveMethodExecutedMatcher) HaveMethodExecutedMatcher {
		m.args = &args
		return m
	}
}

// StoreArgs returns a HaveMethodExecutedOption which stores the arguments
// passed to the method in the addresses provided. A nil value tells the matcher
// to skip the argument at that index.
//
// StoreArgs will cause a panic if:
// - The values provided are not pointers.
// - The mock's arguments are not ConvertibleTo the targets' types.
// - The number of targets does not match the number of arguments in the function.
//   - For variadic methods, a target value must be provided to store the variadic
//     values.  If the value passed in for the variadic values is nil, it will be
//     skipped as normal.  Otherwise, it must be a pointer to a slice capable of
//     storing the variadic arguments' type.
func StoreArgs(targets ...interface{}) HaveMethodExecutedOption {
	return func(m HaveMethodExecutedMatcher) HaveMethodExecutedMatcher {
		m.saveTo = &targets
		return m
	}
}

// Returning returns a HaveMethodExecutedOption which will return the arguments
// on the mock's return channels after the method in question has been called.
//
// Returning will cause a panic if:
//   - The values provided are not ConvertibleTo the mock's return types.
//   - The number of values provided does not match the number of return types in
//     the mock.
//   - For variadic methods, the matcher will use vals[len(nonVariadicArgs(mock)):]
//     as the variadic argument.
func Returning(vals ...interface{}) HaveMethodExecutedOption {
	return func(m HaveMethodExecutedMatcher) HaveMethodExecutedMatcher {
		m.returns = &vals
		return m
	}
}

// HaveMethodExecutedMatcher is a matcher to ensure that a method on a mock was
// executed.
type HaveMethodExecutedMatcher struct {
	MethodName string
	within     time.Duration
	args       *[]interface{}
	saveTo     *[]interface{}
	returns    *[]interface{}

	differ matchers.Differ
}

// HaveMethodExecuted returns a matcher that asserts that the method referenced
// by name was executed.  Options can modify the behavior of the matcher.
//
// HaveMethodExecuted will panic if the mock does not have a method matching
// name.
//
// The HaveMethodExecutedMatcher will panic if any of the options used don't
// match the target mock properly. Check the documentation on the options to get
// more specific information about what would cause the matcher to panic.
func HaveMethodExecuted(name string, opts ...HaveMethodExecutedOption) *HaveMethodExecutedMatcher {
	m := HaveMethodExecutedMatcher{MethodName: name, differ: diff.New()}
	for _, opt := range opts {
		m = opt(m)
	}
	return &m
}

// UseDiffer sets m to use d when showing a diff between actual and expected values.
func (m *HaveMethodExecutedMatcher) UseDiffer(d matchers.Differ) {
	m.differ = d
}

// Match checks the mock value v to see if it has a method matching m.MethodName
// which has been called.
func (m HaveMethodExecutedMatcher) Match(v interface{}) (interface{}, error) {
	if m.differ == nil {
		m.differ = diff.New()
	}
	mv := reflect.ValueOf(v)
	method, exists := mv.Type().MethodByName(m.MethodName)
	if !exists {
		panic(fmt.Errorf("pers: could not find method '%s' on type %T", m.MethodName, v))
	}
	if mv.Kind() == reflect.Ptr {
		mv = mv.Elem()
	}
	if m.returns != nil {
		outField := mv.FieldByName(m.MethodName + "Output")
		defer func() {
			defer func() {
				// The rare double-defer!  This is the only way
				// we can recover from a panic in Return.
				if r := recover(); r != nil {
					panic(fmt.Errorf("pers: HaveMethodExecutedMatcher could not return: %v", r))
				}
			}()
			Return(outField.Interface(), *m.returns...)
		}()
	}
	calledField := mv.FieldByName(m.MethodName + "Called")
	cases := []reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: calledField},
	}
	switch m.within {
	case 0:
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	default:
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(m.within))})
	}

	chosen, _, _ := reflect.Select(cases)
	if chosen == 1 {
		return v, fmt.Errorf("pers: expected method %s to have been called, but it was not", m.MethodName)
	}
	inputField := mv.FieldByName(m.MethodName + "Input")
	if !inputField.IsValid() {
		return v, nil
	}

	if m.saveTo != nil {
		if len(*m.saveTo) != inputField.NumField() {
			msg := fmt.Sprintf("pers: incorrect number of StoreArgs arguments: %v takes %d arguments (got %d)", m.MethodName, inputField.NumField(), len(*m.saveTo))
			if method.Type.IsVariadic() {
				msg = fmt.Sprintf("%s [hint: %v is variadic and you must provide a slice to store all variadic arguments to]", msg, m.MethodName)
			}
			panic(errors.New(msg))
		}
	}

	var calledWith []interface{}
	for i := 0; i < inputField.NumField(); i++ {
		fv, ok := inputField.Field(i).Recv()
		if !ok {
			return v, fmt.Errorf("pers: field %s is closed; cannot perform matches against this mock", inputField.Type().Field(i).Name)
		}
		calledWith = append(calledWith, fv.Interface())

		if m.saveTo != nil {
			tgt := reflect.ValueOf((*m.saveTo)[i]).Elem()
			if !fv.Type().ConvertibleTo(tgt.Type()) {
				panic(fmt.Errorf("pers: wrong type for argument %d in StoreArgs arguments: %v is not convertible to %v", i, fv.Type(), tgt.Type()))
			}
			tgt.Set(fv.Convert(tgt.Type()))
		}
	}
	if m.args == nil {
		return v, nil
	}

	args, err := convertVariadic(method, *m.args)
	if err != nil {
		panic(err)
	}
	if len(args) != len(calledWith) {
		// NOTE: variadic functions were already checked earlier.  This check is for non-variadic functions.
		panic(fmt.Errorf("pers: incorrect number of WithArgs arguments: %v takes exactly %d arguments (got %d)", m.MethodName, len(calledWith), len(args)))
	}

	for i, called := range calledWith {
		arg := args[i]
		switch arg.(type) {
		case nil, any, Matcher:
			continue
		default:
			if method.Type.IsVariadic() && i == len(calledWith)-1 {
				// The variadic value in the expected arguments is always of
				// type []interface{} because we allow values like Any and
				// Matcher, which may not match the element type of the method's
				// variadic slice.
				continue
			}
			v := reflect.ValueOf(arg)
			tgt := reflect.TypeOf(called)
			if !v.Type().ConvertibleTo(tgt) {
				panic(fmt.Errorf("pers: wrong type for argument %d in WithArgs: %v is not convertible to %v", i, v.Type(), tgt))
			}
			args[i] = v.Convert(tgt).Interface()
		}
	}
	matched, diff := m.sliceDiff(reflect.ValueOf(calledWith), reflect.ValueOf(args))
	if matched {
		return v, nil
	}
	const msg = "pers: %s was called with incorrect arguments: %s"
	return v, fmt.Errorf(msg, m.MethodName, diff)
}

func (m HaveMethodExecutedMatcher) sliceDiff(actual, expected reflect.Value) (bool, string) {
	if actual.Len() != expected.Len() {
		return false, m.differ.Diff(fmt.Sprintf("length %d", actual.Len()), fmt.Sprintf("length %d", expected.Len()))
	}
	var diffs []string
	matched := true
	for i := 0; i < actual.Len(); i++ {
		match, diff := m.valueDiff(actual.Index(i), expected.Index(i))
		matched = matched && match
		diffs = append(diffs, diff)
	}
	return matched, fmt.Sprintf("[ %s ]", strings.Join(diffs, ", "))
}

func (m HaveMethodExecutedMatcher) mapDiff(actual, expected reflect.Value) (bool, string) {
	matched := true
	var diffs []string
	for _, k := range expected.MapKeys() {
		eV := expected.MapIndex(k)
		aV := actual.MapIndex(k)
		if !aV.IsValid() {
			matched = false
			diffs = append(diffs, m.differ.Diff("missing key: %v", k.Interface()))
			continue
		}
		match, diff := m.valueDiff(aV, eV)
		matched = matched && match
		diffs = append(diffs, fmt.Sprintf(formatFor(k)+": %s", k.Interface(), diff))
	}
	return matched, fmt.Sprintf("{ %s }", strings.Join(diffs, ", "))
}

func (m HaveMethodExecutedMatcher) valueDiff(actual, expected reflect.Value) (bool, string) {
	for actual.Kind() == reflect.Interface {
		actual = actual.Elem()
	}
	for expected.Kind() == reflect.Interface {
		expected = expected.Elem()
	}
	if !actual.IsValid() || isNil(actual) {
		if !expected.IsValid() || isNil(expected) {
			return true, "<nil>"
		}
	}
	if !expected.IsValid() || isNil(expected) {
		return false, m.differ.Diff(actual.Interface(), nil)
	}

	format := formatFor(actual.Interface())
	actualStr := fmt.Sprintf(format, actual.Interface())
	switch src := expected.Interface().(type) {
	case any:
		return true, actualStr
	case Matcher:
		if dm, ok := src.(expect.DiffMatcher); ok {
			dm.UseDiffer(m.differ)
		}
		_, err := src.Match(actual.Interface())
		if err != nil {
			return false, err.Error()
		}
		return true, actualStr
	default:
		if actual.Kind() != expected.Kind() {
			return false, m.differ.Diff(actual.Interface(), expected.Interface())
		}
		switch actual.Kind() {
		case reflect.Slice, reflect.Array:
			return m.sliceDiff(actual, expected)
		case reflect.Map:
			return m.mapDiff(actual, expected)
		default:
			a, e := actual.Interface(), expected.Interface()
			if !reflect.DeepEqual(a, e) {
				return false, fmt.Sprintf(format, m.differ.Diff(a, e))
			}
			return true, actualStr
		}
	}
}

func convertVariadic(method reflect.Method, args []interface{}) ([]interface{}, error) {
	if !method.Type.IsVariadic() {
		return args, nil
	}

	lastTypeArg := method.Type.NumIn() - 1
	variType := method.Type.In(lastTypeArg).Elem()
	lastArg := lastTypeArg - 1 // lastTypeArg is including the receiver as an argument
	if lastArg > len(args) {
		return nil, fmt.Errorf("pers: incorrect number of WithArgs arguments: %v takes at least %d arguments (got %d)", method.Name, method.Type.NumIn()-2, len(args))
	}
	if lastArg == len(args)-1 && args[lastArg] == VariadicAny {
		return append(args[:lastArg], Any), nil
	}
	variadic := reflect.MakeSlice(reflect.TypeOf([]interface{}(nil)), 0, 0)
	for i := lastArg; i < len(args); i++ {
		arg := args[i]
		argV := reflect.ValueOf(arg)
		if !argV.IsValid() {
			// This was a nil value, but reflect.ValueOf created the invalid value.
			// We need a nil/empty value of the slice's type.
			argV = reflect.Zero(variType)
		}
		switch arg.(type) {
		case any, Matcher:
			variadic = reflect.Append(variadic, argV)
		default:
			if !argV.Type().ConvertibleTo(variType) {
				return nil, fmt.Errorf("pers: variadic argument %[1]v (%[1]T) is not convertible to %[2]v", arg, variType)
			}
			variadic = reflect.Append(variadic, argV.Convert(variType))
		}
	}
	return append(args[:lastArg], variadic.Interface()), nil
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// formatFor returns the format string that should be used for
// the passed in actual type.
func formatFor(actual interface{}) string {
	switch actual.(type) {
	case string:
		return `"%v"`
	default:
		return `%v`

	}
}
