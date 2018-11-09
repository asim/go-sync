// Package task provides an interface for distributed jobs
package task

import (
	"time"
)

// Task represents a distributed task
type Task interface {
	// Run runs a command immediately until completion
	Run(Command) error
	// Schedule defers the execution of a command
	Schedule(Schedule, Command) error
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
