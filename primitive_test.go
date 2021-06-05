package enki

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPrimitive(t *testing.T) {
	p := NewPrimitive()
	require.NotNil(t, p.inner)
	require.NoError(t, p.err)
}

func Test_primitiveBuilder_Err(t *testing.T) {
	p := NewPrimitive()
	p.err = io.EOF
	require.Error(t, p.Err())
}

func Test_primitiveBuilder_NewLine(t *testing.T) {
	p := NewPrimitive()
	p.NewLine()
	require.Greater(t, p.inner.Len(), 0)
}

func Test_primitiveBuilder_String(t *testing.T) {
	p := NewPrimitive()
	p.NewLine()
	require.Equal(t, "\n", p.String())
}
