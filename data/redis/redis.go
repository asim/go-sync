package redis

import (
	"github.com/micro/go-sync/data"
	redis "gopkg.in/redis.v3"
)

type rkv struct {
	Client *redis.Client
}

func (r *rkv) Get(key string) (*data.Item, error) {
	val, err := r.Client.Get(key).Bytes()

	if err != nil && err == redis.Nil {
		return nil, data.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, data.ErrNotFound
	}

	d, err := r.Client.TTL(key).Result()
	if err != nil {
		return nil, err
	}

	return &data.Item{
		Key:        key,
		Value:      val,
		Expiration: d,
	}, nil
}

func (r *rkv) Del(key string) error {
	return r.Client.Del(key).Err()
}

func (r *rkv) Put(item *data.Item) error {
	return r.Client.Set(item.Key, item.Value, item.Expiration).Err()
}

func (r *rkv) List() ([]*data.Item, error) {
	keys, err := r.Client.Keys("*").Result()
	if err != nil {
		return nil, err
	}
	var vals []*data.Item
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

func NewData(opts ...data.Option) data.Data {
	var options data.Options
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
