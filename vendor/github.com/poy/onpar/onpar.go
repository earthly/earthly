//go:generate hel

package onpar

import (
	"errors"
	"fmt"
	"path"
	"testing"
)

type prefs struct {
}

// Opt is an option type to pass to onpar's constructor.
type Opt func(prefs) prefs

type suite interface {
	Run()
	addRunner(runner)
	child() child
}

type child interface {
	addSpecs()
}

// Table is an entry to be used in table tests.
type Table[T, U, V any] struct {
	parent *Onpar[T, U]
	spec   func(U, V)
}

// TableSpec returns a Table type which may be used to declare table tests. The
// spec argument is the test that will be run for each entry in the table.
//
// This is effectively syntactic sugar for looping over table tests and calling
// `parent.Spec` for each entry in the table.
func TableSpec[T, U, V any](parent *Onpar[T, U], spec func(U, V)) Table[T, U, V] {
	return Table[T, U, V]{parent: parent, spec: spec}
}

// Entry adds an entry to t using entry as the value for this table entry.
func (t Table[T, U, V]) Entry(name string, entry V) Table[T, U, V] {
	t.parent.Spec(name, func(v U) {
		t.spec(v, entry)
	})
	return t
}

// FnEntry adds an entry to t that calls setup in order to get its entry value.
// The value from the BeforeEach will be passed to setup, and then both values
// will be passed to the table spec.
func (t Table[T, U, V]) FnEntry(name string, setup func(U) V) Table[T, U, V] {
	t.parent.Spec(name, func(v U) {
		entry := setup(v)
		t.spec(v, entry)
	})
	return t
}

// Onpar stores the state of the specs and groups
type Onpar[T, U any] struct {
	t TestRunner

	path []string

	parent suite

	// level is handled by (*Onpar[T]).Group(), which will adjust this field
	// each time it is called. This is how onpar knows to create nested `t.Run`
	// calls.
	level *level[T, U]

	// canBeforeEach controls which contexts BeforeEach is allowed to take this
	// suite as a parent suite.
	canBeforeEach bool

	// childSuite is assigned by BeforeEach and removed at the end of Group. If
	// BeforeEach is called twice in the same Group (or twice at the top level),
	// this is how it knows to panic.
	//
	// At the end of Group calls, childSuite.addSpecs is called, which will sync the
	// childSuite's specs to the parent.
	childSuite child
	childPath  []string

	runCalled bool
}

// TestRunner matches the methods in *testing.T that the top level onpar
// (returned from New) needs in order to work.
type TestRunner interface {
	Run(name string, fn func(*testing.T)) bool
	Cleanup(func())
}

// New creates a new Onpar suite. The top-level onpar suite must be constructed
// with this. Think `context.Background()`.
//
// It's normal to construct the top-level suite with a BeforeEach by doing the
// following:
//
//	o := BeforeEach(New(t), setupFn)
func New[T TestRunner](t T, opts ...Opt) *Onpar[*testing.T, *testing.T] {
	p := prefs{}
	for _, opt := range opts {
		p = opt(p)
	}
	o := Onpar[*testing.T, *testing.T]{
		t:             t,
		canBeforeEach: true,
		level: &level[*testing.T, *testing.T]{
			before: func(t *testing.T) *testing.T {
				return t
			},
		},
	}
	t.Cleanup(func() {
		if !o.runCalled {
			panic("onpar: Run was never called [hint: missing 'defer o.Run()'?]")
		}
	})
	return &o
}

// Run runs all of o's tests. Typically this will be called in a `defer`
// immediately after o is defined:
//
//	o := onpar.BeforeEach(onpar.New(t), setupFn)
//	defer o.Run()
func (o *Onpar[T, U]) Run() {
	if o.parent == nil {
		o.run(o.t)
		o.runCalled = true
		return
	}
	o.parent.Run()
}

// BeforeEach creates a new child Onpar suite with the requested function as the
// setup function for all specs. It requires a parent Onpar.
//
// The top level Onpar *must* have been constructed with New, otherwise the
// suite will not run.
//
// BeforeEach should be called only once for each level (i.e. each group). It
// will panic if it detects that it is overwriting another BeforeEach call for a
// given level.
func BeforeEach[T, U, V any](parent *Onpar[T, U], setup func(U) V) *Onpar[U, V] {
	if !parent.canBeforeEach {
		panic(fmt.Errorf("onpar: BeforeEach called with invalid parent: parent must either be a top-level suite or be used inside of a `parent.Group()` call"))
	}
	if !parent.correctGroup() {
		panic(fmt.Errorf("onpar: BeforeEach called with invalid parent: parent suite can only be used inside of its group (%v), but the group has exited", path.Join(parent.path...)))
	}
	if parent.child() != nil {
		if len(parent.childPath) == 0 {
			panic(errors.New("onpar: BeforeEach was called more than once at the top level"))
		}
		panic(fmt.Errorf("onpar: BeforeEach was called more than once for group '%s'", path.Join(parent.childPath...)))
	}
	path := parent.path
	if parent.level.name() != "" {
		path = append(parent.path, parent.level.name())
	}
	child := &Onpar[U, V]{
		path:   path,
		parent: parent,
		level: &level[U, V]{
			before: setup,
		},
	}
	parent.childSuite = child
	parent.childPath = child.path
	return child
}

// Spec is a test that runs in parallel with other specs.
func (o *Onpar[T, U]) Spec(name string, f func(U)) {
	if !o.correctGroup() {
		panic(fmt.Errorf("onpar: Spec called on child suite outside of its group (%v)", path.Join(o.path...)))
	}
	spec := concurrentSpec[U]{
		serialSpec: serialSpec[U]{
			specName: name,
			f:        f,
		},
	}
	o.addRunner(spec)
}

// SerialSpec is a test that runs synchronously (i.e. onpar will not call
// `t.Parallel`). While onpar is primarily a parallel testing suite, we
// recognize that sometimes a test just can't be run in parallel. When that is
// the case, use SerialSpec.
func (o *Onpar[T, U]) SerialSpec(name string, f func(U)) {
	if !o.correctGroup() {
		panic(fmt.Errorf("onpar: SerialSpec called on child suite outside of its group (%v)", path.Join(o.path...)))
	}
	spec := serialSpec[U]{
		specName: name,
		f:        f,
	}
	o.addRunner(spec)
}

func (o *Onpar[T, U]) addRunner(r runner) {
	o.level.runners = append(o.level.runners, r)
}

// Group is used to gather and categorize specs. Inside of each group, a new
// child *Onpar may be constructed using BeforeEach.
func (o *Onpar[T, U]) Group(name string, f func()) {
	if !o.correctGroup() {
		panic(fmt.Errorf("onpar: Group called on child suite outside of its group (%v)", path.Join(o.path...)))
	}
	oldLevel := o.level
	o.level = &level[T, U]{
		levelName: name,
	}
	o.canBeforeEach = true
	defer func() {
		o.canBeforeEach = false
		if o.child() != nil {
			o.child().addSpecs()
			o.childSuite = nil
		}
		oldLevel.runners = append(oldLevel.runners,
			&level[U, U]{
				levelName: o.level.name(),
				before: func(v U) U {
					return v
				},
				runners: o.level.runners,
			})
		o.level = oldLevel
	}()

	f()
}

// AfterEach is used to cleanup anything from the specs or BeforeEaches.
// AfterEach may only be called once for each *Onpar value constructed.
func (o *Onpar[T, U]) AfterEach(f func(U)) {
	if !o.correctGroup() {
		panic(fmt.Errorf("onpar: AfterEach called on child suite outside of its group (%v)", path.Join(o.path...)))
	}
	if o.level.after != nil {
		if len(o.childPath) == 0 {
			panic(errors.New("onpar: AfterEach was called more than once at top level"))
		}
		panic(fmt.Errorf("onpar: AfterEach was called more than once for group '%s'", path.Join(o.path...)))
	}
	o.level.after = f
}

func (o *Onpar[T, U]) run(t TestRunner) {
	if o.child() != nil {
		// This happens when New is called before BeforeEach, e.g.:
		//
		//     o := onpar.New()
		//     defer o.Run(t)
		//
		//     b := onpar.BeforeEach(o, setup)
		//
		// Since there's no call to o.Group, the child won't be synced, so we
		// need to do that here.
		o.child().addSpecs()
		o.childSuite = nil
	}
	top, ok := any(o.level).(groupRunner[*testing.T])
	if !ok {
		// This should be impossible - the only place that `run` is called is in
		// `New()`, which is only capable of returning `*Onpar[*testing.T,
		// *testing.T]`.
		var empty T
		panic(fmt.Errorf("onpar: run was called on a child level (type '%T' is not *testing.T)", empty))
	}
	top.runSpecs(t, func() testScope[*testing.T] {
		return baseScope{}
	})
}

func (o *Onpar[T, U]) child() child {
	return o.childSuite
}

func (o *Onpar[T, U]) correctGroup() bool {
	if o.parent == nil {
		return true
	}
	if o.parent.child() == o {
		return true
	}
	return false
}

// addSpecs is called by parent Group() calls to tell o to add its specs to its
// parent.
func (o *Onpar[T, U]) addSpecs() {
	o.parent.addRunner(o.level)
}

type testScope[T any] interface {
	before(*testing.T) T
	after()
}

type baseScope struct {
}

func (s baseScope) before(t *testing.T) *testing.T {
	return t
}

func (s baseScope) after() {}

type runner interface {
	name() string
}

type groupRunner[T any] interface {
	runner
	runSpecs(t TestRunner, scope func() testScope[T])
}

type specRunner[T any] interface {
	runner
	run(t *testing.T, scope func() testScope[T])
}

type concurrentSpec[T any] struct {
	serialSpec[T]
}

func (s concurrentSpec[T]) run(t *testing.T, scope func() testScope[T]) {
	t.Parallel()

	s.serialSpec.run(t, scope)
}

type serialSpec[T any] struct {
	specName string
	f        func(T)
}

func (s serialSpec[T]) name() string {
	return s.specName
}

func (s serialSpec[T]) run(t *testing.T, scope func() testScope[T]) {
	sc := scope()
	v := sc.before(t)
	s.f(v)
	sc.after()
}

type levelScope[T, U any] struct {
	val          U
	parentBefore func(*testing.T) T
	childBefore  func(T) U
	childAfter   func(U)
	parentAfter  func()
}

func (s *levelScope[T, U]) before(t *testing.T) U {
	parentVal := s.parentBefore(t)
	s.val = s.childBefore(parentVal)
	return s.val
}

func (s *levelScope[T, U]) after() {
	if s.childAfter != nil {
		s.childAfter(s.val)
	}
	if s.parentAfter != nil {
		s.parentAfter()
	}
}

type level[T, U any] struct {
	levelName string
	before    func(T) U
	after     func(U)
	runners   []runner
}

func (l *level[T, U]) name() string {
	return l.levelName
}

func (l *level[T, U]) runSpecs(t TestRunner, scope func() testScope[T]) {
	for _, r := range l.runners {
		childScope := func() testScope[U] {
			parentScope := scope()
			return &levelScope[T, U]{
				parentBefore: parentScope.before,
				childBefore:  l.before,
				childAfter:   l.after,
				parentAfter:  parentScope.after,
			}
		}
		switch r := r.(type) {
		case groupRunner[U]:
			if r.name() == "" {
				// If the name is empty, running the group as a sub-group would
				// result in ugly output. Just run the test function at this level
				// instead.
				r.runSpecs(t, childScope)
				return
			}
			t.Run(r.name(), func(t *testing.T) {
				r.runSpecs(t, childScope)
			})
		case specRunner[U]:
			t.Run(r.name(), func(t *testing.T) {
				r.run(t, childScope)
			})
		default:
			panic(fmt.Errorf("onpar: spec runner type [%T] is not supported", r))
		}
	}
}
