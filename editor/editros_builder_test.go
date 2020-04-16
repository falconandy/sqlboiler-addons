package editor

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditorsBuilder_Build(t *testing.T) {
	var buf bytes.Buffer
	buf.Reset()
	b := NewBuilder(filepath.Join("..", "_testdata"))
	b.FileExt = ".go.txt"
	err := b.Build(&buf)
	require.Nil(t, err)

	expected, err := ioutil.ReadFile(filepath.Join("..", "_testdata", "expected", "boil_editors.go.txt"))
	require.Nil(t, err)

	assert.Equal(t, string(expected), buf.String())
}
