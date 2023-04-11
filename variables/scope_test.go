package variables_test

import (
	"testing"

	"github.com/earthly/earthly/variables"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
)

func TestScope(topT *testing.T) {
	type testCtx struct {
		t      *testing.T
		expect expect.Expectation
		scope  *variables.Scope
	}

	o := onpar.BeforeEach(onpar.New(topT), func(t *testing.T) testCtx {
		return testCtx{
			t:      t,
			expect: expect.New(t),
			scope:  variables.NewScope(),
		}
	})
	defer o.Run()

	o.Spec("it returns false for unset variables", func(tc testCtx) {
		_, ok := tc.scope.Get("foo")
		tc.expect(ok).To(beFalse())
	})

	o.Spec("NoOverride prevents Add from overriding an existing value", func(tc testCtx) {
		tc.scope.Add("foo", "bar")
		tc.scope.Add("foo", "baz", variables.WithActive(), variables.NoOverride())

		v, ok := tc.scope.Get("foo")
		tc.expect(ok).To(beTrue())
		tc.expect(v).To(equal("bar"))

		_, ok = tc.scope.Get("foo", variables.WithActive())
		tc.expect(ok).To(beFalse())
	})

	o.Spec("it returns a sorted list of names", func(tc testCtx) {
		tc.scope.Add("a", "", variables.WithActive())
		tc.scope.Add("z", "", variables.WithActive())
		tc.scope.Add("e", "")
		tc.scope.Add("b", "", variables.WithActive())

		inactive := tc.scope.Sorted()
		tc.expect(inactive).To(equal([]string{"a", "b", "e", "z"}))
		active := tc.scope.Sorted(variables.WithActive())
		tc.expect(active).To(equal([]string{"a", "b", "z"}))
	})

	for _, tt := range []struct {
		testName    string
		useOpts     []variables.ScopeOpt
		failGetOpts []variables.ScopeOpt
		name        string
		value       string
	}{
		{
			testName: "it stores inactive values",
			failGetOpts: []variables.ScopeOpt{
				variables.WithActive(),
			},
			name:  "foo",
			value: "bar",
		},
		{
			testName: "it stores active values",
			useOpts:  []variables.ScopeOpt{variables.WithActive()},
			name:     "bar",
			value:    "baz",
		},
		{
			testName: "it stores active env variables",
			useOpts: []variables.ScopeOpt{
				variables.WithActive(),
			},
			name:  "bacon",
			value: "eggs",
		},
	} {
		tt := tt
		o.Spec(tt.testName, func(tc testCtx) {
			ok := tc.scope.Add(tt.name, tt.value)
			tc.expect(ok).To(beTrue())
			for _, opt := range tt.useOpts {
				_, ok := tc.scope.Get(tt.name, opt)
				tc.expect(ok).To(beFalse())
				ok = tc.scope.Add(tt.name, tt.value, opt, variables.NoOverride())
				tc.expect(ok).To(beFalse())
				ok = tc.scope.Add(tt.name, tt.value, opt)
				tc.expect(ok).To(beTrue())
			}

			value, ok := tc.scope.Get(tt.name)
			tc.expect(ok).To(beTrue())
			tc.expect(value).To(equal(tt.value))
			for _, opt := range tt.useOpts {
				value, ok := tc.scope.Get(tt.name, opt)
				tc.expect(ok).To(beTrue())
				tc.expect(value).To(equal(tt.value))

				m := tc.scope.Map(opt)
				value, ok = m[tt.name]
				tc.expect(ok).To(beTrue())
				tc.expect(value).To(equal(tt.value))
			}

			for _, opt := range tt.failGetOpts {
				_, ok := tc.scope.Get(tt.name, opt)
				tc.expect(ok).To(beFalse())

				m := tc.scope.Map(opt)
				_, ok = m[tt.name]
				tc.expect(ok).To(beFalse())
			}

			clone := tc.scope.Clone()
			value, ok = clone.Get(tt.name)
			tc.expect(ok).To(beTrue())
			tc.expect(value).To(equal(tt.value))
			for _, opt := range tt.useOpts {
				value, ok = clone.Get(tt.name, opt)
				tc.expect(ok).To(beTrue())
				tc.expect(value).To(equal(tt.value))
			}

			tc.scope.Remove(tt.name)
			tc.scope.Add(tt.name, tt.value)
			for _, opt := range tt.useOpts {
				_, ok := tc.scope.Get(tt.name, opt)
				tc.expect(ok).To(beFalse())
			}
		})
	}

	o.Group("CombineScopes", func() {
		o.Spec("it prefers left values", func(tc testCtx) {
			tc.scope.Add("a", "b")

			other := variables.NewScope()
			other.Add("a", "c")

			c := variables.CombineScopes(tc.scope, other)
			v, ok := c.Get("a")
			tc.expect(ok).To(beTrue())
			tc.expect(v).To(equal("b"))
		})

		o.Spec("it prefers active to inactive values", func(tc testCtx) {
			tc.scope.Add("active", "b")

			other := variables.NewScope()
			other.Add("active", "d", variables.WithActive())

			c := variables.CombineScopes(tc.scope, other)
			env, ok := c.Get("active")
			tc.expect(ok).To(beTrue())
			tc.expect(env).To(equal("d"))
		})
	})
}
