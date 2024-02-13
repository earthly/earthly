package authprovider_test

import (
	"time"

	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/poy/onpar/matchers"
)

const (
	timeout     = time.Second
	mockTimeout = 5 * time.Second
)

var (
	beTrue       = matchers.BeTrue
	beFalse      = matchers.BeFalse
	not          = matchers.Not
	haveOccurred = matchers.HaveOccurred
	equal        = matchers.Equal
	matchRegexp  = matchers.MatchRegexp
	beNil        = matchers.BeNil

	haveMethodExecuted = pers.HaveMethodExecuted
	within             = pers.Within
	withArgs           = pers.WithArgs
	returning          = pers.Returning
)
