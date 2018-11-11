// Package task provides an interface for distributed jobs
package task

import (
	"time"
)

// Task represents a distributed task
type Task interface {
	// Run runs a command immediately until completion
	Run(Command) error
	// Status provides status of last execution
	Status() string
}

type Command func() error

// Schedule represents a time or interval at which a task should run
type Schedule struct {
	// When to start the schedule. Zero time means immediately
	Time time.Time
	// Non zero interval dictates an ongoing schedule
	Interval time.Duration
}

type Options struct {
	Pool int
}

type Option func(o *Options)

func (s Schedule) Run() <-chan time.Time {
	d := s.Time.Sub(time.Now())

	ch := make(chan time.Time, 1)

	go func() {
		// wait for start time
		<-time.After(d)

		// zero interval
		if s.Interval == time.Duration(0) {
			ch <- time.Now()
			close(ch)
			return
		}

		// start ticker
		for t := range time.Tick(s.Interval) {
			ch <- t
		}
	}()

	return ch
}
