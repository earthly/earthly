package cloud_test

import (
	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/poy/onpar/matchers"
)

var (
	haveOccurred = matchers.HaveOccurred
	not          = matchers.Not
	equal        = matchers.Equal

	haveMethodExecuted = pers.HaveMethodExecuted
	withArgs           = pers.WithArgs
	storeArgs          = pers.StoreArgs
	within             = pers.Within
	returning          = pers.Returning
)
