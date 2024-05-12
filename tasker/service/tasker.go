package service

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/teran/backupnizza/models"
)

var (
	ErrAlreadyExists = errors.New("task already exists")
	ErrNotFound      = errors.New("not found")
)

type Tasker interface {
	Add(task models.Task) error
	GetByName(name string) (models.Runner, error)
}

type tasker struct {
	tasks map[string]models.Runner
	mutex *sync.Mutex
}

func New() Tasker {
	return &tasker{
		tasks: make(map[string]models.Runner),
		mutex: &sync.Mutex{},
	}
}

func (t *tasker) Add(task models.Task) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, ok := t.tasks[task.Name]; ok {
		return errors.Wrap(ErrAlreadyExists, task.Name)
	}

	t.tasks[task.Name] = task.Runner
	return nil
}

func (t *tasker) GetByName(name string) (models.Runner, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	task, ok := t.tasks[name]
	if !ok {
		return nil, errors.Wrap(ErrNotFound, name)
	}

	return task, nil
}
