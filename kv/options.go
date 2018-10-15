package kv

type Options struct {
	Addrs []string
}

func Addrs(a ...string) Option {
	return func(o *Options) {
		o.Addrs = a
	}
}
