// Package consul is a consul implementation of kv
package consul

import (
	"fmt"
	"net"

	"github.com/hashicorp/consul/api"
	"github.com/micro/go-sync/store"
)

type ckv struct {
	client *api.Client
}

func (c *ckv) Get(key string) (*store.Item, error) {
	keyval, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}

	if keyval == nil {
		return nil, store.ErrNotFound
	}

	return &store.Item{
		Key:   keyval.Key,
		Value: keyval.Value,
	}, nil
}

func (c *ckv) Del(key string) error {
	_, err := c.client.KV().Delete(key, nil)
	return err
}

func (c *ckv) Put(item *store.Item) error {
	_, err := c.client.KV().Put(&api.KVPair{
		Key:   item.Key,
		Value: item.Value,
	}, nil)
	return err
}

func (c *ckv) List() ([]*store.Item, error) {
	keyval, _, err := c.client.KV().List("/", nil)
	if err != nil {
		return nil, err
	}
	if keyval == nil {
		return nil, store.ErrNotFound
	}
	var vals []*store.Item
	for _, keyv := range keyval {
		vals = append(vals, &store.Item{
			Key:   keyv.Key,
			Value: keyv.Value,
		})
	}
	return vals, nil
}

func (c *ckv) String() string {
	return "consul"
}

func NewStore(opts ...store.Option) store.Store {
	var options store.Options
	for _, o := range opts {
		o(&options)
	}

	config := api.DefaultConfig()

	// set host
	// config.Host something
	// check if there are any addrs
	if len(options.Nodes) > 0 {
		addr, port, err := net.SplitHostPort(options.Nodes[0])
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "8500"
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		} else if err == nil {
			config.Address = fmt.Sprintf("%s:%s", addr, port)
		}
	}

	client, _ := api.NewClient(config)

	return &ckv{
		client: client,
	}
}
