// Package sync is a distributed synchronization framework
package sync

import (
	"github.com/micro/go-sync/data"
	"github.com/micro/go-sync/leader"
	"github.com/micro/go-sync/lock"
	"github.com/micro/go-sync/task"
	"github.com/micro/go-sync/time"
)

// DB provides synchronized access to key-value storage.
// It uses the data interface and lock interface to
// provide a consistent storage mechanism.
type DB interface {
	// Load value with given key
	Load(key, val interface{}) error
	// Store value with given key
	Store(key, val interface{}) error
	// Delete value with given key
	Delete(key interface{}) error
	// Range over all key/vals. Value changes are saved
	Range(func(key, val interface{}) error) error
}

// Cron is a distributed scheduler using leader election
// and distributed task runners. It uses the leader and
// task interfaces.
type Cron interface {
	Schedule(task.Schedule, task.Command) error
}

type Options struct {
	Leader leader.Leader
	Lock   lock.Lock
	Data   data.Data
	Task   task.Task
	Time   time.Time
}

type Option func(o *Options)
