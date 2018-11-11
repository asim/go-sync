package sync

import (
	"github.com/micro/go-sync/task"
)

type syncCron struct {
	opts Options
}

func (c *syncCron) Schedule(s task.Schedule, t task.Command) error {
	go func() {
		// run the scheduler
		tc := s.Run()

		// execute the task
		for _ = range tc {
			c.opts.Task.Run(t)
		}
	}()

	return nil
}
