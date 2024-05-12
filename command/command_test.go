package command

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/teran/backupnizza/command/environment"
)

func TestRun(t *testing.T) {
	r := require.New(t)

	stdout := &bytes.Buffer{}

	c := NewCommand("testdata/pong", WithStdout(stdout))

	err := c.Run(context.TODO())
	r.NoError(err)
	r.Equal("pong\n", stdout.String())
}

func TestRunWithEnvironment(t *testing.T) {
	r := require.New(t)

	stdout := &bytes.Buffer{}

	c := NewCommand("testdata/print_env_var", WithStdout(stdout), WithEnvironment(map[string]environment.Value{
		"TEST_VARIABLE": func() (string, error) { return "some value", nil },
	}))

	err := c.Run(context.TODO())
	r.NoError(err)
	r.Equal("some value\n", stdout.String())
}

func TestRunWithArgs(t *testing.T) {
	r := require.New(t)

	stdout := &bytes.Buffer{}

	c := NewCommand("testdata/args_echo", WithStdout(stdout), WithArgs("-f=value1", "arg1"))

	err := c.Run(context.TODO())
	r.NoError(err)
	r.Equal("-f=value1 arg1\n", stdout.String())
}

func TestRunWithError(t *testing.T) {
	r := require.New(t)

	stderr := &bytes.Buffer{}

	c := NewCommand("testdata/inexecutable", WithStderr(stderr))

	err := c.Run(context.TODO())
	r.Error(err)
	r.Equal("error running command: fork/exec testdata/inexecutable: permission denied", err.Error())
}
