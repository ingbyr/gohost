package host

import "sync"

var (
	instance *service
	once     sync.Once
)

func Service() *service {
	once.Do(func() {
		instance = newService()
	})
	return instance
}

type service struct {
}

func newService() *service {
	return &service{}
}
