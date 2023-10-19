package diff

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/poy/onpar/diff/str"
)

var DefaultStrDiffs = []StringDiffAlgorithm{str.NewCharDiff()}

const DefaultTimeout = time.Second

// StringDiffAlgorithm is a type which can generate diffs between two strings.
type StringDiffAlgorithm interface {
	// Diffs takes a context.Context to know when to stop, returning a channel
	// of diffs. Each new diff returned on this channel should have a lower
	// cost than the previous one.
	//
	// If a higher cost diff is returned after a lower cost diff, it will be
	// discarded.
	//
	// Once ctx.Done() is closed, diffs will not be read off of the returned
	// channel - it's up to the algorithm to perform select statements to avoid
	// deadlocking.
	Diffs(ctx context.Context, actual, expected []rune) <-chan str.Diff
}

// WithStringAlgos picks the algorithms that will be used to diff strings. We
// will always use a "dumb" algorithm for the base case that simply returns the
// full string as either equal or different, but extra algorithms can be
// provided to get more useful diffs for larger strings. For example,
// StringCharDiff gets diffs of characters (good for catching misspellings), and
// StringLineDiff gets diffs of lines (good for large multiline output).
//
// The default is DefaultStrDiffs.
//
// If called without any arguments, only the "dumb" algorithm will be used.
func WithStringAlgos(algos ...StringDiffAlgorithm) Opt {
	return func(d Differ) Differ {
		d.stringAlgos = algos
		return d
	}
}

// WithTimeout sets a timeout for diffing. Normally, diffs will be refined until
// the "cost" of the diff is as low as possible. If diffing is taking too long,
// the best diff that has been loaded will be returned.
//
// The default is DefaultTimeout.
//
// If no diff at all has been generated yet, we will still wait until the first
// diff is generated before returning, but no refinement will be done.
func WithTimeout(timeout time.Duration) Opt {
	return func(d Differ) Differ {
		d.timeout = timeout
		return d
	}
}

// Opt is an option type that can be passed to New.
//
// Most of the time, you'll want to use at least one
// of Actual or Expected, to differentiate the two
// in your output.
type Opt func(Differ) Differ

// WithFormat returns an Opt that wraps up differences
// using a format string.  The format should contain
// one '%s' to add the difference string in.
func WithFormat(format string) Opt {
	return func(d Differ) Differ {
		d.wrappers = append(d.wrappers, func(v string) string {
			return fmt.Sprintf(format, v)
		})
		return d
	}
}

// Sprinter is any type which can print a string.
type Sprinter interface {
	Sprint(...any) string
}

// WithSprinter returns an Opt that wraps up differences
// using a Sprinter.
func WithSprinter(s Sprinter) Opt {
	return func(d Differ) Differ {
		d.wrappers = append(d.wrappers, func(v string) string {
			return s.Sprint(v)
		})
		return d
	}
}

func applyOpts(o *Differ, opts ...Opt) {
	for _, opt := range opts {
		*o = opt(*o)
	}
}

// Actual returns an Opt that only applies formatting to the actual value.
// Non-formatting options (e.g. different diffing algorithms) will have no
// effect.
func Actual(opts ...Opt) Opt {
	return func(d Differ) Differ {
		if d.actual == nil {
			d.actual = &Differ{}
		}
		applyOpts(d.actual, opts...)
		return d
	}
}

// Expected returns an Opt that only applies formatting to the expected value.
// Non-formatting options (e.g. different diffing algorithms) will have no
// effect.
func Expected(opts ...Opt) Opt {
	return func(d Differ) Differ {
		if d.expected == nil {
			d.expected = &Differ{}
		}
		applyOpts(d.expected, opts...)
		return d
	}
}

// Differ is a type that can diff values.  It keeps its own
// diffing style.
type Differ struct {
	wrappers []func(string) string

	actual   *Differ
	expected *Differ

	timeout     time.Duration
	stringAlgos []StringDiffAlgorithm
}

// New creates a Differ, using the passed in opts to manipulate
// its diffing behavior and output.
//
// By default, we wrap mismatched text in angle brackets and separate them with
// "!=". Example:
//
//     matching text >actual!=expected< more matching text
//
// opts will be applied to the text in the order they
// are passed in, so you can do things like color a value
// and then wrap the colored text up in custom formatting.
//
// See the examples on the different Opt types for more
// detail.
func New(opts ...Opt) *Differ {
	d := Differ{
		timeout:     DefaultTimeout,
		stringAlgos: DefaultStrDiffs,
	}
	for _, o := range opts {
		d = o(d)
	}
	if d.needsDefaultFmt() {
		d = WithFormat(">%s<")(d)
		d = Actual(WithFormat("%s!="))(d)
	}
	return &d
}

func (d *Differ) needsDefaultFmt() bool {
	return len(d.wrappers) == 0 &&
		d.actual == nil &&
		d.expected == nil
}

// format is used to format a string using the wrapper functions.
func (d Differ) format(v string) string {
	for _, w := range d.wrappers {
		v = w(v)
	}
	return v
}

// Diff takes two values and returns a string showing a
// diff of them.
func (d *Differ) Diff(actual, expected any) string {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	return d.diff(ctx, reflect.ValueOf(actual), reflect.ValueOf(expected))
}

func (d *Differ) genDiff(format string, actual, expected any) string {
	afmt := fmt.Sprintf(format, actual)
	if d.actual != nil {
		afmt = d.actual.format(afmt)
	}
	efmt := fmt.Sprintf(format, expected)
	if d.expected != nil {
		efmt = d.expected.format(efmt)
	}
	return d.format(afmt + efmt)
}

func (d *Differ) diff(ctx context.Context, av, ev reflect.Value) string {
	if !av.IsValid() {
		if !ev.IsValid() {
			return "<nil>"
		}
		if ev.Kind() == reflect.Ptr {
			return d.diff(ctx, av, ev.Elem())
		}
		return d.genDiff("%v", "<nil>", ev.Interface())
	}
	if !ev.IsValid() {
		if av.Kind() == reflect.Ptr {
			return d.diff(ctx, av.Elem(), ev)
		}
		return d.genDiff("%v", av.Interface(), "<nil>")
	}

	if av.Kind() != ev.Kind() {
		return d.genDiff("%s", av.Type(), ev.Type())
	}

	if av.CanInterface() {
		switch av.Interface().(type) {
		case []rune, []byte, string:
			// TODO: we probably want to (eventually) run this concurrently. As
			// is, a struct with two complicated strings in two separate fields
			// that both need diffs will probably get a pretty good diff for the
			// first field and just the baseline diff for the second one.
			return d.strDiff(ctx, av, ev)
		}
	}

	switch av.Kind() {
	case reflect.Ptr, reflect.Interface:
		return d.diff(ctx, av.Elem(), ev.Elem())
	case reflect.Slice, reflect.Array, reflect.String:
		if av.Len() != ev.Len() {
			// TODO: do a more thorough diff of values
			return d.genDiff(fmt.Sprintf("%s(len %%d)", av.Type()), av.Len(), ev.Len())
		}
		var elems []string
		for i := 0; i < av.Len(); i++ {
			elems = append(elems, d.diff(ctx, av.Index(i), ev.Index(i)))
		}
		return "[ " + strings.Join(elems, ", ") + " ]"
	case reflect.Map:
		var parts []string
		for _, kv := range ev.MapKeys() {
			k := kv.Interface()
			bmv := ev.MapIndex(kv)
			amv := av.MapIndex(kv)
			if !amv.IsValid() {
				parts = append(parts, d.genDiff("%s", fmt.Sprintf("missing key %v", k), fmt.Sprintf("%v: %v", k, bmv.Interface())))
				continue
			}
			parts = append(parts, fmt.Sprintf("%v: %s", k, d.diff(ctx, amv, bmv)))
		}
		for _, kv := range av.MapKeys() {
			// We've already compared all keys that exist in both maps; now we're
			// just looking for keys that only exist in a.
			if !ev.MapIndex(kv).IsValid() {
				k := kv.Interface()
				parts = append(parts, d.genDiff("%s", fmt.Sprintf("extra key %v: %v", k, av.MapIndex(kv).Interface()), fmt.Sprintf("%v: nil", k)))
				continue
			}
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case reflect.Struct:
		if av.Type().Name() != ev.Type().Name() {
			return d.genDiff("%s", av.Type().Name(), ev.Type().Name()) + "(mismatched types)"
		}
		var parts []string
		for i := 0; i < ev.NumField(); i++ {
			f := ev.Type().Field(i)
			if f.PkgPath != "" && !f.Anonymous {
				// unexported
				continue
			}
			name := f.Name
			bfv := ev.Field(i)
			afv := av.Field(i)
			parts = append(parts, fmt.Sprintf("%s: %s", name, d.diff(ctx, afv, bfv)))
		}
		return fmt.Sprintf("%s{%s}", av.Type(), strings.Join(parts, ", "))
	default:
		if av.Type().Comparable() {
			a, b := av.Interface(), ev.Interface()
			if a != b {
				return d.genDiff("%#v", a, b)
			}
			return fmt.Sprintf("%#v", a)
		}
		return d.format(fmt.Sprintf("UNSUPPORTED: could not compare values of type %s", av.Type()))
	}
}

// strDiff helps us generate a diff between two strings. It uses the baseStrAlgo
// first to get a baseline, then uses results from the other algorithms set in
// d.stringAlgos to get the lowest cost possible before returning.
func (d *Differ) strDiff(ctx context.Context, av, ev reflect.Value) string {
	runeTyp := reflect.TypeOf([]rune(nil))
	actual := av.Convert(runeTyp).Interface().([]rune)
	expected := ev.Convert(runeTyp).Interface().([]rune)

	var wg sync.WaitGroup
	results := make(chan str.Diff)
	for _, algo := range d.stringAlgos {
		algoCh := algo.Diffs(ctx, actual, expected)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case diff, ok := <-algoCh:
					if !ok {
						return
					}
					results <- diff
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		// All of our algorithms are sending results to the same channel. Once
		// they are all done (either from the timeout or from exhausting all
		// options), we need to close the results channel to let the main logic
		// know that everything's done - we don't want to continue waiting if
		// all algorithms have exhausted their possible diffs.
		//
		// Since we know that the results channel will be closed once the
		// context times out, there's no reason to select on the results
		// channel.
		defer close(results)
		wg.Wait()
	}()

	best := baseStringDiff(actual, expected)
	for diff := range results {
		if diff.Cost() >= best.Cost() {
			continue
		}
		best = diff
	}

	var out string
	for _, section := range best.Sections() {
		if section.Type == str.TypeMatch {
			out += string(section.Actual)
			continue
		}
		out += d.genDiff("%s", string(section.Actual), string(section.Expected))
	}
	return out
}
