package memcached

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	mc "github.com/bradfitz/gomemcache/memcache"
	"github.com/micro/go-sync/data"
)

type mkv struct {
	Server *mc.ServerList
	Client *mc.Client
}

func (m *mkv) Get(key string) (*data.Item, error) {
	keyval, err := m.Client.Get(key)
	if err != nil && err == mc.ErrCacheMiss {
		return nil, data.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if keyval == nil {
		return nil, data.ErrNotFound
	}

	return &data.Item{
		Key:        keyval.Key,
		Value:      keyval.Value,
		Expiration: time.Second * time.Duration(keyval.Expiration),
	}, nil
}

func (m *mkv) Del(key string) error {
	return m.Client.Delete(key)
}

func (m *mkv) Put(item *data.Item) error {
	return m.Client.Set(&mc.Item{
		Key:        item.Key,
		Value:      item.Value,
		Expiration: int32(item.Expiration.Seconds()),
	})
}

func (m *mkv) List() ([]*data.Item, error) {
	// stats
	// cachedump
	// get keys

	var keys []string

	//data := make(map[string]string)
	if err := m.Server.Each(func(c net.Addr) error {
		cc, err := net.Dial("tcp", c.String())
		if err != nil {
			return err
		}
		defer cc.Close()

		b := bufio.NewReadWriter(bufio.NewReader(cc), bufio.NewWriter(cc))

		// get items
		if _, err := fmt.Fprintf(b, "stats items\r\n"); err != nil {
			return err
		}

		b.Flush()

		v, err := b.ReadSlice('\n')
		if err != nil {
			return err
		}

		parts := bytes.Split(v, []byte("\n"))
		if len(parts) < 1 {
			return nil
		}
		vals := strings.Split(string(parts[0]), ":")
		items := vals[1]

		// drain
		for {
			buf, err := b.ReadSlice('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if strings.HasPrefix(string(buf), "END") {
				break
			}
		}

		b.Writer.Reset(cc)
		b.Reader.Reset(cc)

		if _, err := fmt.Fprintf(b, "lru_crawler metadump %s\r\n", items); err != nil {
			return err
		}
		b.Flush()

		for {
			v, err := b.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if strings.HasPrefix(v, "END") {
				break
			}
			key := strings.Split(v, " ")[0]
			keys = append(keys, strings.TrimPrefix(key, "key="))
		}

		return nil
	}); err != nil {
		return nil, err
	}

	var vals []*data.Item

	// concurrent op
	ch := make(chan *data.Item, len(keys))

	for _, k := range keys {
		go func(key string) {
			i, _ := m.Get(key)
			ch <- i
		}(k)
	}

	for i := 0; i < len(keys); i++ {
		item := <-ch

		if item == nil {
			continue
		}

		vals = append(vals, item)
	}

	close(ch)

	return vals, nil
}

func (m *mkv) String() string {
	return "memcached"
}

func NewData(opts ...data.Option) data.Data {
	var options data.Options
	for _, o := range opts {
		o(&options)
	}

	if len(options.Nodes) == 0 {
		options.Nodes = []string{"127.0.0.1:11211"}
	}

	ss := new(mc.ServerList)
	ss.SetServers(options.Nodes...)

	return &mkv{
		Server: ss,
		Client: mc.New(options.Nodes...),
	}
}
