package sync

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/micro/go-sync/kv"
	ckv "github.com/micro/go-sync/kv/consul"
	lock "github.com/micro/go-sync/lock/consul"
)

type syncMap struct {
	opts Options
}

func ekey(k interface{}) string {
	b, _ := json.Marshal(k)
	return base64.StdEncoding.EncodeToString(b)
}

func (m *syncMap) Load(key, val interface{}) error {
	if key == nil {
		return fmt.Errorf("key is nil")
	}

	kstr := ekey(key)

	// lock
	if err := m.opts.Lock.Acquire(kstr); err != nil {
		return err
	}
	defer m.opts.Lock.Release(kstr)

	// get key
	kval, err := m.opts.KV.Get(kstr)
	if err != nil {
		return err
	}

	// decode value
	return json.Unmarshal(kval.Value, val)
}

func (m *syncMap) Store(key, val interface{}) error {
	if key == nil {
		return fmt.Errorf("key is nil")
	}

	kstr := ekey(key)

	// lock
	if err := m.opts.Lock.Acquire(kstr); err != nil {
		return err
	}
	defer m.opts.Lock.Release(kstr)

	// encode value
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	// set key
	return m.opts.KV.Put(&kv.Item{
		Key:   kstr,
		Value: b,
	})
}

func (m *syncMap) Delete(key interface{}) error {
	if key == nil {
		return fmt.Errorf("key is nil")
	}

	kstr := ekey(key)

	// lock
	if err := m.opts.Lock.Acquire(kstr); err != nil {
		return err
	}
	defer m.opts.Lock.Release(kstr)
	return m.opts.KV.Del(kstr)
}

func (m *syncMap) Range(fn func(key, val interface{}) error) error {
	keyvals, err := m.opts.KV.List()
	if err != nil {
		return err
	}

	for _, keyval := range keyvals {
		// lock
		if err := m.opts.Lock.Acquire(keyval.Key); err != nil {
			return err
		}
		// unlock
		defer m.opts.Lock.Release(keyval.Key)

		// unmarshal value
		var val interface{}

		if len(keyval.Value) > 0 && keyval.Value[0] == '{' {
			if err := json.Unmarshal(keyval.Value, &val); err != nil {
				return err
			}
		} else {
			val = keyval.Value
		}

		// exec func
		if err := fn(keyval.Key, val); err != nil {
			return err
		}

		// save val
		b, err := json.Marshal(val)
		if err != nil {
			return err
		}

		// no save
		if i := bytes.Compare(keyval.Value, b); i == 0 {
			return nil
		}

		// set key
		if err := m.opts.KV.Put(&kv.Item{
			Key:   keyval.Key,
			Value: b,
		}); err != nil {
			return err
		}
	}

	return nil
}

func NewMap(opts ...Option) Map {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	if options.Lock == nil {
		options.Lock = lock.NewLock()
	}

	if options.KV == nil {
		options.KV = ckv.NewKV()
	}

	return &syncMap{
		opts: options,
	}
}
