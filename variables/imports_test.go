package variables_test

import (
	"github.com/poy/onpar/matchers"
	"github.com/stretchr/testify/require"
)

var (
	beTrue       = matchers.BeTrue
	beFalse      = matchers.BeFalse
	equal        = matchers.Equal
	not          = matchers.Not
	haveOccurred = matchers.HaveOccurred
)

var (
	NoError = require.NoError
	Error   = require.Error
	Equal   = require.Equal
)
