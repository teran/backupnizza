package command

import (
	"context"
	"encoding/json"
	"os"
	"runtime"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/teran/backupnizza/command/environment"
	"github.com/teran/backupnizza/command/environment/providers/secretbox"
	"github.com/teran/backupnizza/command/environment/providers/static"
	"github.com/teran/backupnizza/models"
	secretboxService "github.com/teran/secretbox/service"
)

type Group struct {
	commands       []models.Runner
	executionLogic ExecutionLogic
}

func NewGroupFromJSON(svc secretboxService.Service, sandboxCliPath string, in json.RawMessage) (*Group, error) {
	cfg := GroupOptions{}
	err := json.Unmarshal(in, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling configuration")
	}

	err = cfg.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "error validating configuration")
	}

	cmds := []models.Runner{}
	for _, cmd := range cfg.Subtasks {
		env := make(map[string]environment.Value)
		for _, v := range cmd.Environment {
			if v.ValueFrom != nil {
				switch *v.ValueFrom {
				case "secretbox_value":
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
							cmd.Name, opts.Name,
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
							cmd.Name, opts.Name,
						)
					}

					env[v.Name] = secretbox.NewSecretboxCommandVariable(svc, sandboxCliPath, opts.Socket, opts.Name)
				default:
					return nil, errors.Errorf("unexpected value from value: `%s`", *v.ValueFrom)
				}
			} else {
				env[v.Name] = static.NewStaticVariable(*v.Value)
			}
		}

		log.WithFields(log.Fields{
			"component": "command runner",
			"task_name": cmd.Name,
		}).Info("registering subtask")

		cmds = append(cmds, NewCommand(
			cmd.Binary,
			WithArgs(cmd.Arguments...),
			WithCopyCurrentProcessEnvironment(),
			WithEnvironment(env),
			WithStderr(os.Stderr),
			WithStdout(os.Stdout),
		))
	}

	return &Group{
		executionLogic: cfg.ExecutionLogic,
		commands:       cmds,
	}, nil
}

func (c *Group) Run(ctx context.Context) error {
	switch c.executionLogic {
	case ExecutionLogicSequentialAll:
		return c.runSequentialAll(ctx)
	case ExecutionLogicSequentialDontCare:
		return c.runSequentialDontCare(ctx)
	case ExecutionLogicParallelAll:
		return c.runParallelAll(ctx)
	}

	return errors.Errorf("unknown execution logic: %s", string(c.executionLogic))
}

func (c *Group) runSequentialAll(ctx context.Context) error {
	for _, cmd := range c.commands {
		err := cmd.Run(ctx)
		if err != nil {
			return errors.Wrap(err, "error running chained task")
		}
	}
	return nil
}

func (c *Group) runSequentialDontCare(ctx context.Context) error {
	for i, cmd := range c.commands {
		err := cmd.Run(ctx)
		if err != nil {
			log.Debugf("error running subcommand#%d: %s", i, err)
		}
	}

	return nil
}

func (c *Group) runParallelAll(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	for _, cmd := range c.commands {
		g.Go(func(ctx context.Context, cmd models.Runner) func() error {
			return func() error {
				return cmd.Run(ctx)
			}
		}(ctx, cmd))
	}

	return g.Wait()
}
