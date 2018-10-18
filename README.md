# Go Sync [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/go-sync?status.svg)](https://godoc.org/github.com/micro/go-sync)

Go Sync is a synchronization framework for distributed systems.

## Overview

Distributed systems by their very nature are decoupled and independent. In most cases they must honour 2 out of 3 letters of the CAP theorem 
e.g Availability and Partitional tolerance but sacrificing consistency. In the case of microservices we often offload this concern to 
an external database or eventing system. Go Sync provides a framework for synchronization which can be used in the application by the developer.

We offer three primitives; Lock, Leader and KV

- Lock -  distributed locking 
- Leader - leadership election
- KV - key value storage

## Getting Started

- [Locking](#locking) - exclusive resource access
- [Leadership](#leadership) - single leader group coordination
- [Key-Value](#key-value) - simple distributed data storage

## Locking

The Lock interface provides distributed locking. Multiple instances attempting to lock the same id will block until available.

```go
import "github.com/micro/go-sync/lock/consul"

l := consul.NewLock()

// acquire lock
err := lock.Acquire("id")
// handle err

// release lock
err = lock.Release("id")
// handle err
```

## Leadership

Leader provides leadership election. Useful where one node needs to coordinate some action.

```go
import (
	"github.com/micro/go-sync/leader"
	"github.com/micro/go-sync/leader/consul"
)

l := consul.NewLeader(
	leader.Group("name"),
)

// elect leader
e, err := l.Elect("id")
// handle err


// operate while leader
revoked := e.Revoked()

for {
	select {
	case <-revoked:
		// re-elect
		e.Elect("id")
	default:
		// leader operation
	}
}

// resign leadership
e.Resign() 
```

## Key-Value

KV provides a simple interface for distributed key-value stores.

```go
import (
	"github.com/micro/go-sync/kv"
	"github.com/micro/go-sync/kv/consul"
)

keyval := consul.NewKV()

err := keyval.Put(&kv.Item{
	Key: "foo",
	Value: []byte(`bar`),
})
// handle err

err = keyval.Get("foo")
// handle err

err = keyval.Delete("foo")
```
