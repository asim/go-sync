package lock

import (
	"time"
)

// Addrs sets the addresses the underlying lock implementation
func Addrs(a ...string) Option {
	return func(o *Options) {
		o.Addrs = append(o.Addrs, a...)
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
