package command

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/teran/backupnizza/command/environment"
	"github.com/teran/backupnizza/command/environment/providers/secretbox"
	"github.com/teran/backupnizza/command/environment/providers/static"
	secretboxService "github.com/teran/secretbox/service"
)

type Command struct {
	binaryPath  string
	args        []string
	environment map[string]environment.Value
	stdout      io.Writer
	stdin       io.Reader
	stderr      io.Writer
}

func NewCommandFromJSON(svc secretboxService.Service, cliPath string, in json.RawMessage) (*Command, error) {
	cfg := CommandConfig{}
	err := json.Unmarshal(in, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling configuration")
	}

	env := make(map[string]environment.Value)
	for _, v := range cfg.Environment {
		if v.ValueFrom != nil {
			switch *v.ValueFrom {
			case "secretbox_value":
				type secretboxOptions struct {
					Name string `json:"name"`
				}

				opts := secretboxOptions{}
				err := json.Unmarshal(v.Options, &opts)
				if err != nil {
					return nil, err
				}

				ok, err := svc.IsSecretRegistered(context.TODO(), opts.Name)
				if err != nil {
					return nil, errors.Wrap(err, "error checking if the secret is registered")
				}

				if !ok {
					return nil, errors.Errorf(
						"secret is not registered but mentioned in configuration: cmd_name=`%s` secret_name=`%s`",
						cfg.Name, opts.Name,
					)
				}

				env[v.Name] = secretbox.NewSecretboxVariable(svc, opts.Name)
			case "secretbox_command":
				type secretboxOptions struct {
					Name   string `json:"name"`
					Socket string `json:"socket"`
				}

				opts := secretboxOptions{}
				err := json.Unmarshal(v.Options, &opts)
				if err != nil {
					return nil, err
				}

				ok, err := svc.IsSecretRegistered(context.TODO(), opts.Name)
				if err != nil {
					return nil, errors.Wrap(err, "error checking if the secret is registered")
				}

				if !ok {
					return nil, errors.Errorf(
						"secret is not registered but mentioned in configuration: cmd_name=`%s` secret_name=`%s`",
						cfg.Name, opts.Name,
					)
				}

				env[v.Name] = secretbox.NewSecretboxCommandVariable(svc, cliPath, opts.Socket, opts.Name)
			default:
				return nil, errors.Errorf("unexpected value from value: `%s`", *v.ValueFrom)
			}
		} else {
			env[v.Name] = static.NewStaticVariable(*v.Value)
		}
	}

	return NewCommand(
		cfg.Binary,
		WithArgs(cfg.Arguments...),
		WithCopyCurrentProcessEnvironment(),
		WithEnvironment(env),
		WithStderr(os.Stderr),
		WithStdout(os.Stdout),
	), nil
}

func NewCommand(binaryPath string, opts ...Option) *Command {
	c := &Command{
		binaryPath:  binaryPath,
		environment: make(map[string]environment.Value),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (cmd *Command) Run(ctx context.Context) error {
	c := exec.CommandContext(ctx, cmd.binaryPath, cmd.args...)
	env, err := mapMapToSlice(cmd.environment)
	if err != nil {
		return errors.Wrap(err, "error compiling environment variables set")
	}

	c.Env = env

	log.WithFields(log.Fields{
		"component":   "command_runner",
		"binary":      cmd.binaryPath,
		"args":        cmd.args,
		"environment": env,
	}).Trace("running command")

	if cmd.stdin != nil {
		c.Stdin = cmd.stdin
	}

	if cmd.stdout != nil {
		c.Stdout = cmd.stdout
	}

	if cmd.stderr != nil {
		c.Stderr = cmd.stderr
	}

	return errors.Wrap(c.Run(), "error running command")
}
