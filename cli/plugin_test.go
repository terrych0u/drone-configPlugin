package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingBranch(t *testing.T) {
	plugin := Plugin{}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, missingBranch, err.Error())
}
