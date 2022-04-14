package rainhari

type Service interface {
	Name() string
	Init(...Options)
	Run() error
}

type Option func(*Options)

func NewService(opts ...Option) Service {
	return newService(opts...)
}
