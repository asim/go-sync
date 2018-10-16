package lock

import (
	"time"
)

// Nodes sets the addresses the underlying lock implementation
func Nodes(a ...string) Option {
	return func(o *Options) {
		o.Nodes = a
	}
}

// TTL sets the lock ttl
func TTL(t time.Duration) AcquireOption {
	return func(o *AcquireOptions) {
		o.TTL = t
	}
}

// Wait sets the wait time
func Wait(t time.Duration) AcquireOption {
	return func(o *AcquireOptions) {
		o.Wait = t
	}
}
