package etcd

import (
	"context"
	"log"

	"github.com/micro/go-sync/kv"
	client "go.etcd.io/etcd/clientv3"
)

type ekv struct {
	kv client.KV
}

func (e *ekv) Get(key string) (*kv.Item, error) {
	keyval, err := e.kv.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	if keyval == nil || len(keyval.Kvs) == 0 {
		return nil, kv.ErrNotFound
	}

	return &kv.Item{
		Key:   string(keyval.Kvs[0].Key),
		Value: keyval.Kvs[0].Value,
	}, nil
}

func (e *ekv) Del(key string) error {
	_, err := e.kv.Delete(context.Background(), key)
	return err
}

func (e *ekv) Put(item *kv.Item) error {
	_, err := e.kv.Put(context.Background(), item.Key, string(item.Value))
	return err
}

func (e *ekv) String() string {
	return "etcd"
}

func NewKV(opts ...kv.Option) kv.KV {
	var options kv.Options
	for _, o := range opts {
		o(&options)
	}

	var endpoints []string

	for _, addr := range options.Addrs {
		if len(addr) > 0 {
			endpoints = append(endpoints, addr)
		}
	}

	if len(endpoints) == 0 {
		endpoints = []string{"http://127.0.0.1:2379"}
	}

	// TODO: parse addresses
	c, err := client.New(client.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &ekv{
		kv: client.NewKV(c),
	}
}
