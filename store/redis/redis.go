package redis

import (
	"github.com/micro/go-sync/store"
	redis "gopkg.in/redis.v3"
)

type rkv struct {
	Client *redis.Client
}

func (r *rkv) Get(key string) (*store.Item, error) {
	val, err := r.Client.Get(key).Bytes()

	if err != nil && err == redis.Nil {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, store.ErrNotFound
	}

	d, err := r.Client.TTL(key).Result()
	if err != nil {
		return nil, err
	}

	return &store.Item{
		Key:        key,
		Value:      val,
		Expiration: d,
	}, nil
}

func (r *rkv) Del(key string) error {
	return r.Client.Del(key).Err()
}

func (r *rkv) Put(item *store.Item) error {
	return r.Client.Set(item.Key, item.Value, item.Expiration).Err()
}

func (r *rkv) List() ([]*store.Item, error) {
	keys, err := r.Client.Keys("*").Result()
	if err != nil {
		return nil, err
	}
	var vals []*store.Item
	for _, k := range keys {
		i, err := r.Get(k)
		if err != nil {
			return nil, err
		}
		vals = append(vals, i)
	}
	return vals, nil
}

func (r *rkv) String() string {
	return "redis"
}

func NewStore(opts ...store.Option) store.Store {
	var options store.Options
	for _, o := range opts {
		o(&options)
	}

	if len(options.Nodes) == 0 {
		options.Nodes = []string{"127.0.0.1:6379"}
	}

	return &rkv{
		Client: redis.NewClient(&redis.Options{
			Addr:     options.Nodes[0],
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}
