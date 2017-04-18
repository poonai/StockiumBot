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
		).Info("SERVICE")
	}(time.Now())
	return l.s.echo(req)
}

func (l logger) sendSuggestion(id string, msg string) {
	defer func(begin time.Time) {
		l.log.WithFields(
			logrus.Fields{
				"start_At":  begin.String(),
				"end_at":    time.Since(begin),
				"service":   "sendSuggestion",
				"message":   msg,
				"sender_id": id,
			},
		).Info("SERVICE")
	}(time.Now())
	l.s.sendSuggestion(id, msg)
}

func (l logger) sendFinancialData(id string, companyID string) {
	defer func(begin time.Time) {
		l.log.WithFields(
			logrus.Fields{
				"start_at":    begin.String(),
				"end_at":      time.Since(begin),
				"service":     "sendFinancialData",
				"sender_id":   id,
				"company_url": companyID,
			},
		).Info("SERVICE")
	}(time.Now())
	l.s.sendFinancialData(id, companyID)
}

func (l logger) addToWatchlist(senderID string, companyURL string) {
	defer func(begin time.Time) {
		l.log.WithFields(logrus.Fields{
			"star_at":     begin.String(),
			"end_at":      time.Since(begin),
			"service":     "addToWatchlist",
			"sender_id":   senderID,
			"company_url": companyURL,
		}).Info("SERVICE")

	}(time.Now())
	l.s.addToWatchlist(senderID, companyURL)
}
