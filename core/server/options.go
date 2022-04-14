package server

type Options struct {
	Name    string
	Address string
	Version string
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Address(a string) Option {
	return func(o * Options) {
		o.Address = a
	}
}

func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}
