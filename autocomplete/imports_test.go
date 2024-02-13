package autocomplete

import (
	"github.com/poy/onpar/matchers"
	"github.com/stretchr/testify/assert"
)

var (
	equal        = matchers.Equal
	not          = matchers.Not
	haveOccurred = matchers.HaveOccurred
)

var (
	NoError = assert.NoError
	Equal   = assert.Equal
	Nil     = assert.Nil
	True    = assert.True
	False   = assert.False
	Error   = assert.Error
)
