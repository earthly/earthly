package logstreamer_test

import (
	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/poy/onpar/matchers"
)

var (
	not          = matchers.Not
	haveOccurred = matchers.HaveOccurred

	haveMethodExecuted = pers.HaveMethodExecuted
	within             = pers.Within
	storeArgs          = pers.StoreArgs
)
