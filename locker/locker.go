package locker

import (
	"sync"
)

type Locker interface {
	Lock()
	Unlock()
}

type locker struct {
	mutex *sync.Mutex
}

func New() Locker {
	return &locker{
		mutex: &sync.Mutex{},
	}
}

func (l *locker) Lock() {
	l.mutex.Lock()
}

func (l *locker) Unlock() {
	l.mutex.Unlock()
}
