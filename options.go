package rainhari

import "github.com/agus7fauzi/rainhari/core/server"

type Options struct {
	Server server.Server
}

func newOptions(opts ...Options) Options {
	opt := Options{
		Server: server,
	}
}

func Name(n string) Option {
	return func(o *Options) {
		o.Server.Init(server.Name(n))
	}
}
