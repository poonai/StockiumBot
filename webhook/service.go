package webhook

import (
	"github.com/sch00lb0y/StockiumBot/fb"
)

// Service haing interface of service
type Service interface {
	echo(id string, msg string) string
}

type service struct {
}

// NewService sd
func NewService() service {
	return service{}
}

func (s service) echo(id string, msg string) string {
	go fb.Send(id, msg)
	return msg
}
