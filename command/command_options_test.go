package command

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/teran/backupnizza/command/environment"
)

func TestWithArgs(t *testing.T) {
	r := require.New(t)

	c := NewCommand("blah")

	WithArgs("-f=blah", "test_arg")(c)

	r.Equal([]string{"-f=blah", "test_arg"}, c.args)
}

func TestWithEnvironment(t *testing.T) {
	r := require.New(t)

	c := NewCommand("blah")

	WithEnvironment(map[string]environment.Value{
		"key1": func() (string, error) { return "value1", nil },
		"key2": func() (string, error) { return "value2", nil },
	})(c)

	res, err := mapMapToSlice(c.environment)
	r.NoError(err)
	r.ElementsMatch([]string{
		"key1=value1",
		"key2=value2",
	}, res)
}

func TestWithCopyCurrentProcessEnvironment(t *testing.T) {
	r := require.New(t)

	c := NewCommand("blah")

	WithCopyCurrentProcessEnvironment()(c)

	res, err := mapMapToSlice(c.environment)
	r.NoError(err)
	r.ElementsMatch(os.Environ(), res)
}

func TestWithStdin(t *testing.T) {
	r := require.New(t)

	c := NewCommand("blah")

	r.Nil(c.stdin)

	WithStdin(&bytes.Buffer{})(c)

	r.NotNil(c.stdin)
}

func TestWithStdout(t *testing.T) {
	r := require.New(t)

	c := NewCommand("blah")

	r.Nil(c.stdout)

	WithStdout(&bytes.Buffer{})(c)

	r.NotNil(c.stdout)
}

func TestWithStderr(t *testing.T) {
	r := require.New(t)

	c := NewCommand("blah")

	r.Nil(c.stderr)

	WithStderr(&bytes.Buffer{})(c)

	r.NotNil(c.stderr)
}
