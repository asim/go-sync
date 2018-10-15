# KV [![GoDoc](https://godoc.org/github.com/micro/go-os?status.svg)](https://godoc.org/github.com/micro/go-os/kv)
 
Provides a high level abstraction for key-value stores.

## Interface

```go
type KV interface {
	Close() error
	Get(key string) (*Item, error)
	Del(key string) error
	Put(item *Item) error
	String() string
}

type Item struct {
	Key        string
	Value      []byte
	Expiration time.Duration
}

func NewKV(opts ...Option) KV {
	return newOS(opts...)
}
```

## Supported Backends

- Broker
- [KV Service](https://github.com/micro/kv-srv)
- [Memcached](https://github.com/micro/go-plugins/tree/master/kv/memcached)
- [Redis](https://github.com/micro/go-plugins/tree/master/kv/redis)

## Usage

KV provides a simple Key-Value abstraction. The default implementation a form of broadcast for announcement and in memory key-value via RPC. 
It's embedded within the service as the KV handler and has less than a second startup time. Ring participation is done via the 
go-micro/broker and a gossip topic which can be changed with the variable kv.GossipTopic. Alternatively you can use redis or memcached.

Here's an example using the in-memory distributed store.

```go
package main

import (
        "fmt"
        "time"

        "github.com/micro/go-micro"
        "github.com/micro/go-os/kv"
)

func main() {
        // Create a new service. 
        service := micro.NewService()

        // Create KV instance
        keyval := kv.NewKV(
                kv.Client(service.Client()),
                kv.Server(service.Server()),
        )
        defer keyval.Close()

	// write some keys and try get them back
        go func() {
                for i := 0; i < 10; i++ {
                        time.Sleep(time.Second)
                        key := fmt.Sprintf("greeting-%d", i)
                        err := keyval.Put(&kv.Item{
                                Key:        key,
                                Value:      []byte(`hello`),
                                Expiration: time.Second * 20,
                        })
                        if err != nil {
                                fmt.Println("put err", err)
                                continue
                        }

                        item, err := keyval.Get(key)
                        if err != nil {
                                fmt.Println("get err", err)
                                continue
                        }
                        fmt.Println("got", item)
                }
        }()

        service.Run()
}
```
