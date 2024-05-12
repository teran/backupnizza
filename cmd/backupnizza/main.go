package main

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	grpcServer "google.golang.org/grpc"

	"github.com/teran/backupnizza/command"
	"github.com/teran/backupnizza/config"
	"github.com/teran/backupnizza/locker"
	"github.com/teran/backupnizza/logger"
	"github.com/teran/backupnizza/models"
	"github.com/teran/backupnizza/scheduler"
	onepassword "github.com/teran/go-onepassword-cli"
	"github.com/teran/secretbox/presenter/grpc"
	secretsRepository "github.com/teran/secretbox/repository/secrets/memory"
	tokensRepository "github.com/teran/secretbox/repository/tokens/memory"
	"github.com/teran/secretbox/service"
)

type spec struct {
	Config string `arg:"-c,env:CONFIG" default:"~/.config/backup/config.yaml"`
}

func main() {
	ctx := context.TODO()

	var rcfg spec
	arg.MustParse(&rcfg)

	configPath, err := expandPath(rcfg.Config)
	if err != nil {
		panic(err)
	}

	cfg, err := config.NewFromFile(configPath)
	if err != nil {
		panic(err)
	}

	logCfg, err := logger.NewConfigFromJSON(cfg.Log)
	if err != nil {
		panic(err)
	}

	if err := logger.SetupFromConfig(logCfg); err != nil {
		panic(err)
	}

	secretsRepo := secretsRepository.New()
	tokensRepo := tokensRepository.New()
	svc := service.New(secretsRepo, tokensRepo)

	for _, secret := range cfg.Secretbox.Secrets {
		secretValue := ""
		switch secret.Source {
		case "onepassword":
			type opOptions struct {
				Label string `json:"label"`
				Kind  string `json:"kind"`
			}

			opOpts := opOptions{}
			err = json.Unmarshal(secret.Options, &opOpts)
			if err != nil {
				panic(errors.Wrap(err, "error decoding onepassword options"))
			}

			var kind onepassword.Kind
			switch opOpts.Kind {
			case "credential":
				kind = onepassword.KindCredential
			default:
				kind = onepassword.KindPassword
			}

			secretValue, err = onepassword.New().GetByLabel(ctx, kind, opOpts.Label)
			if err != nil {
				panic(errors.Wrap(err, "error retrieving secret from onepassword"))
			}
		default:
			panic(errors.Errorf("unexpected secret source: `%s`", secret.Source))
		}

		err = svc.CreateSecret(ctx, secret.Name, secretValue)
		if err != nil {
			panic(errors.Wrap(err, "error registering secret in secretbox"))
		}
	}

	taskLocker := locker.New()
	tasks := []scheduler.Task{}

	log.WithFields(log.Fields{
		"name":     "tokens repository cleanup task",
		"schedule": cfg.Secretbox.TokensGCSchedule,
	}).Info("adding to schedule")
	tasks = append(tasks, scheduler.Task{
		Schedule: "* * * * *",
		Task: models.Task{
			Name:   "tokens repository cleanup",
			Runner: scheduler.NewRunnerAdapter(tokensRepo.Cleanup),
		},
	})

	for _, task := range cfg.Tasks {
		log.WithFields(log.Fields{
			"name":     task.Name,
			"schedule": task.Schedule,
		}).Info("adding to schedule")

		var runner models.Runner

		switch task.Kind {
		case "command":
			runner, err = command.NewCommandFromJSON(svc, cfg.Secretbox.CLIPath, task.Options)
			if err != nil {
				panic(err)
			}
		case "group":
			runner, err = command.NewGroupFromJSON(svc, cfg.Secretbox.CLIPath, task.Options)
			if err != nil {
				panic(err)
			}
		default:
			panic(errors.Errorf("unexpected task kind: `%s`", task.Kind))
		}

		tasks = append(tasks, scheduler.Task{
			Schedule: task.Schedule,
			Task: models.Task{
				Name:   task.Name,
				Runner: locker.NewRunnerWithLocks(taskLocker, runner),
			},
		})
	}

	s, err := scheduler.New(time.UTC, tasks...)
	if err != nil {
		panic(err)
	}

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		err := s.Run(ctx)
		return errors.Wrap(err, "error running scheduler")
	})

	presenter := grpc.New(svc)

	gs := grpcServer.NewServer()
	presenter.Register(gs)

	if _, err := os.Stat(cfg.Secretbox.Socket); !os.IsNotExist(err) {
		if err := os.Remove(cfg.Secretbox.Socket); err != nil {
			panic(err)
		}
	}

	listener, err := net.Listen("unix", cfg.Secretbox.Socket)
	if err != nil {
		panic(err)
	}

	err = os.Chmod(cfg.Secretbox.Socket, 0o700)
	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"socket": cfg.Secretbox.Socket,
	}).Infof("starting secretbox server")

	wg.Go(func() error {
		err := gs.Serve(listener)
		return errors.Wrap(err, "error running secretbox API server")
	})

	if err := wg.Wait(); err != nil {
		panic(err)
	}
}

func expandPath(s string) (string, error) {
	if strings.HasPrefix(s, "~/") {
		u, err := user.Current()
		if err != nil {
			return "", errors.Wrap(err, "error getting current user")
		}

		return path.Join(u.HomeDir, strings.TrimLeft(s, "~")), nil
	}
	return s, nil
}
