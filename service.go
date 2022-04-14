package rainhari

type service struct {
	opts Options
}

func newService(opts ...Option) Service {
	service := new(service)
	options := newOptions(opts...)
}