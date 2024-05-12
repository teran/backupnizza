package locker

import (
	"context"

	"github.com/teran/backupnizza/models"
)

var _ models.Runner = (*runner)(nil)

type runner struct {
	locker Locker
	runner models.Runner
}

func NewRunnerWithLocks(l Locker, r models.Runner) models.Runner {
	return &runner{
		locker: l,
		runner: r,
	}
}

func (r *runner) Run(ctx context.Context) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	return r.runner.Run(ctx)
}
