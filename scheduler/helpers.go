package scheduler

import "context"

type RunnerAdapter struct {
	fn func(ctx context.Context) error
}

func NewRunnerAdapter(fn func(ctx context.Context) error) *RunnerAdapter {
	return &RunnerAdapter{
		fn: fn,
	}
}

func (ra *RunnerAdapter) Run(ctx context.Context) error {
	return ra.fn(ctx)
}
