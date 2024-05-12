package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/teran/backupnizza/models"
)

func newGoCronTask(runner models.Runner) func(string, gocron.Job) {
	return func(taskName string, job gocron.Job) {
		log.WithFields(log.Fields{
			"component": component,
			"task_name": taskName,
		}).Info("running task")

		start := time.Now()
		if err := runner.Run(job.Context()); err != nil {
			log.WithFields(log.Fields{
				"component": component,
				"task_name": taskName,
			}).Warnf("error received while running task: %s", err.Error())
		}
		end := time.Now()

		log.WithFields(log.Fields{
			"component":  component,
			"task_name":  taskName,
			"time_taken": end.Sub(start),
		}).Info("task run completed")
	}
}
