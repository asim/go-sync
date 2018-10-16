package redis

import (
	"github.com/micro/go-sync/kv"
	redis "gopkg.in/redis.v3"
)

type rkv struct {
	Client *redis.Client
}

func (r *rkv) Close() error {
	return nil
}

func (r *rkv) Get(key string) (*kv.Item, error) {
	val, err := r.Client.Get(key).Bytes()

	if err != nil && err == redis.Nil {
		return nil, kv.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, kv.ErrNotFound
	}

	d, err := r.Client.TTL(key).Result()
	if err != nil {
		return nil, err
	}

	return &kv.Item{
		Key:        key,
		Value:      val,
		Expiration: d,
	}, nil
}

func (r *rkv) Del(key string) error {
	return r.Client.Del(key).Err()
}

func (r *rkv) Put(item *kv.Item) error {
	return r.Client.Set(item.Key, item.Value, item.Expiration).Err()
}

func (r *rkv) String() string {
	return "redis"
}

func NewKV(opts ...kv.Option) kv.KV {
	var options kv.Options
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
