package service

import (
	"github.com/stretchr/testify/mock"
	"github.com/teran/backupnizza/models"
)

var _ Tasker = (*Mock)(nil)

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Add(task models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *Mock) GetByName(name string) (models.Runner, error) {
	args := m.Called(name)
	return args.Get(0).(models.Runner), args.Error(1)
}
