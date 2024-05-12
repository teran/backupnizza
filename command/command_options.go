package command

import (
	"io"
	"os"
	"strings"

	"github.com/teran/backupnizza/command/environment"
)

func WithArgs(args ...string) Option {
	return func(c *Command) {
		c.args = args
	}
}

func WithEnvironment(env map[string]environment.Value) Option {
	return func(c *Command) {
		c.environment = env
	}
}

func WithCopyCurrentProcessEnvironment() Option {
	return func(c *Command) {
		for _, v := range os.Environ() {
			parts := strings.SplitN(v, "=", 2)
			if len(parts) != 2 {
				continue
			}

			c.environment[parts[0]] = func(value string) environment.Value {
				return func() (string, error) {
					return value, nil
				}
			}(parts[1])
		}
	}
}

func WithStdin(rd io.Reader) Option {
	return func(c *Command) {
		c.stdin = rd
	}
}

func WithStdout(wr io.Writer) Option {
	return func(c *Command) {
		c.stdout = wr
	}
}

func WithStderr(wr io.Writer) Option {
	return func(c *Command) {
		c.stderr = wr
	}
}
