package ollama

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	o := New(nil)
	require.NotNil(t, o)
}
