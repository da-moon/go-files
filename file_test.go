package files_test

import (
	"testing"

	files "github.com/da-moon/go-files"
	"github.com/stretchr/testify/assert"
)

func TestIsTemporaryFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "", expected: false},
		{input: ".", expected: false},
		{input: ".tmp", expected: false},
		{input: "tmp.123", expected: false},
		{input: ".tmp.123.xx", expected: false},
		{input: "asdf.sdfds.tmp.dfd", expected: false},
		{input: "dfd.sdfds.dfds.1232", expected: false},
		{input: ".tmp.1", expected: true},
		{input: "asdf.dff.tmp.123", expected: true},
	}

	for _, tt := range tests {
		actual := files.IsTemporaryFileName(tt.input)
		assert.Equal(t, tt.expected, actual)
	}
}
