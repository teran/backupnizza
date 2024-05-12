package static

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStaticVariable(t *testing.T) {
	r := require.New(t)

	fn := NewStaticVariable("test_string")
	v, err := fn()
	r.NoError(err)
	r.Equal("test_string", v)
}
