package sync

import (
	"github.com/micro/go-sync/task"
)

type syncCron struct {
	opts Options
}

func (c *syncCron) Schedule(s task.Schedule, t task.Command) error {
	return c.opts.Task.Schedule(s, t)
}
