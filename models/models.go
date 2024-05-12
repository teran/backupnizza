package models

import "context"

type Runner interface {
	Run(ctx context.Context) error
}

type Task struct {
	Name   string
	Runner Runner
}
