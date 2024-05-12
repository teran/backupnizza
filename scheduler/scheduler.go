package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/teran/backupnizza/models"
)

type Task struct {
	Schedule string
	Task     models.Task
}

type Scheduler interface {
	Run(ctx context.Context) error
}

type scheduler struct {
	cron *gocron.Scheduler
}

func New(loc *time.Location, tasks ...Task) (Scheduler, error) {
	cron := gocron.NewScheduler(loc)

	for _, task := range tasks {
		_, err := cron.
			Cron(task.Schedule).
			DoWithJobDetails(newGoCronTask(task.Task.Runner), task.Task.Name)
		if err != nil {
			return nil, err
		}
	}

	return &scheduler{
		cron: cron,
	}, nil
}

func (s *scheduler) Run(ctx context.Context) error {
	s.cron.StartBlocking()
	return nil
}
