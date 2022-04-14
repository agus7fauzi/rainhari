package rainhari

type Service interface {
	Init()
	Run()
}

type Option func(*Options)

func NewService() Service {

}
