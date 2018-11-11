package sync

import (
	"fmt"
	"math"
	"time"

	"github.com/micro/go-log"
	"github.com/micro/go-sync/lock/consul"
	"github.com/micro/go-sync/task"
	"github.com/micro/go-sync/task/local"
)

type syncCron struct {
	opts Options
}

func backoff(attempts int) time.Duration {
	if attempts == 0 {
		return time.Duration(0)
	}
	return time.Duration(math.Pow(10, float64(attempts))) * time.Millisecond
}

func (c *syncCron) Schedule(s task.Schedule, t task.Command) error {
	id := fmt.Sprintf("%s-%s", s.String(), t.String())

	go func() {
		// run the scheduler
		tc := s.Run()

		var i int

		for {
			// leader election
			e, err := c.opts.Leader.Elect(id)
			if err != nil {
				log.Logf("[cron] leader election error: %v", err)
				time.Sleep(backoff(i))
				i++
				continue
			}

			i = 0
			r := e.Revoked()

			// execute the task
			for {
				select {
				// schedule tick
				case <-tc:
					c.opts.Task.Run(t)
				// leader revoked
				case <-r:
					break
				}
			}
		}
	}()

	return nil
}

func NewCron(opts ...Option) Cron {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	if options.Lock == nil {
		options.Lock = consul.NewLock()
	}

	if options.Task == nil {
		options.Task = local.NewTask()
	}

	return &syncCron{
		opts: options,
	}
}
