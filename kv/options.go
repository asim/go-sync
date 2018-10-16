package kv

type Options struct {
	Nodes []string
}

func Nodes(a ...string) Option {
	return func(o *Options) {
		o.Nodes = a
	}
}
