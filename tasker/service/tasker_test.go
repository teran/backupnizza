package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teran/backupnizza/models"
)

func TestTasker(t *testing.T) {
	r := require.New(t)

	// Create new tasker
	tt := New()

	// Add new task
	runner1 := &testRunner{}
	err := tt.Add(models.Task{
		Name:   "test task",
		Runner: runner1,
	})
	r.NoError(err)

	// Add the same task twice
	runner2 := &testRunner{}
	err = tt.Add(models.Task{
		Name:   "test task",
		Runner: runner2,
	})
	r.Error(err)
	errors.Is(err, ErrAlreadyExists)

	// Get & run existent one
	retrRunner, err := tt.GetByName("test task")
	r.NoError(err)

	err = retrRunner.Run(context.Background())
	r.NoError(err)

	// Check if first runner was ran and second one is not
	r.True(runner1.isRan)
	r.False(runner2.isRan)

	// Get not existent one
	_, err = tt.GetByName("not-existent-task")
	r.Error(err)
	errors.Is(err, ErrNotFound)

	// Just some greybox checks since we simply can
	r.Len(tt.(*tasker).tasks, 1)
}

type testRunner struct {
	isRan bool
}

func (tr *testRunner) Run(ctx context.Context) error {
	tr.isRan = true
	return nil
}
