// Package lock provides distributed locking
package lock

import (
	"time"
)

// Lock is a distributed locking interfac
type Lock interface {
	Acquire(id string, opts ...AcquireOption) error
	Release(id string) error
	String() string
}

type Options struct {
	Addrs []string
}

type AcquireOptions struct {
	TTL  time.Duration
	Wait time.Duration
}

type Option func(o *Options)
type AcquireOption func(o *AcquireOptions)
