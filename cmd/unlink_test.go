package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Unlink(t *testing.T) {

	actual := new(bytes.Buffer)
	unlinkCmd.SetOut(actual)
	unlinkCmd.SetErr(actual)
	unlinkCmd.SetArgs([]string{"unlink"})
	unlinkCmd.Execute()

	expected := ""

	assert.Equal(t, expected, actual.String())
}
