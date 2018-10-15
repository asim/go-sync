package memcached

import (
	"time"

	mc "github.com/bradfitz/gomemcache/memcache"
	"github.com/micro/go-sync/kv"
)

type mkv struct {
	Client *mc.Client
}

func (m *mkv) Get(key string) (*kv.Item, error) {
	keyval, err := m.Client.Get(key)
	if err != nil && err == mc.ErrCacheMiss {
		return nil, kv.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if keyval == nil {
		return nil, kv.ErrNotFound
	}

	return &kv.Item{
		Key:        keyval.Key,
		Value:      keyval.Value,
		Expiration: time.Second * time.Duration(keyval.Expiration),
	}, nil
}

func (m *mkv) Del(key string) error {
	return m.Client.Delete(key)
}

func (m *mkv) Put(item *kv.Item) error {
	return m.Client.Set(&mc.Item{
		Key:        item.Key,
		Value:      item.Value,
		Expiration: int32(item.Expiration.Seconds()),
	})
}

func (m *mkv) String() string {
	return "memcached"
}

func NewKV(opts ...kv.Option) kv.KV {
	var options kv.Options
	for _, o := range opts {
		o(&options)
	}

	if len(options.Addrs) == 0 {
		options.Addrs = []string{"127.0.0.1:11211"}
	}

	return &mkv{
		Client: mc.New(options.Addrs...),
	}
}
