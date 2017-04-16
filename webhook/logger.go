package webhook

import (
	"time"

	"github.com/Sirupsen/logrus"
)

type logger struct {
	s   Service
	log *logrus.Logger
}

// NewLogger is used to return logger Service
func NewLogger(log *logrus.Logger, s Service) Service {
	return logger{s, log}
}

func (l logger) echo(req request) string {
	defer func(begin time.Time) {
		l.log.WithFields(
			logrus.Fields{
				"start_At": begin.String(),
				"end_at":   time.Since(begin).String(),
				"service":  "echo",
				"message":  req.Entry[0].Messaging[0].Message.Text,
			},
		).Info("Echo service")
	}(time.Now())
	return l.s.echo(req)
}
