package authprovider_test

import (
	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/poy/onpar/matchers"
)

var (
	beTrue       = matchers.BeTrue
	not          = matchers.Not
	haveOccurred = matchers.HaveOccurred
	equal        = matchers.Equal
	beClosed     = matchers.BeClosed
	matchRegexp  = matchers.MatchRegexp

	haveMethodExecuted = pers.HaveMethodExecuted
	within             = pers.Within
	withArgs           = pers.WithArgs
	returning          = pers.Returning
)
