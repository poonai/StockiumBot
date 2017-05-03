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

func (l logger) echo(senderID string, message string) string {
	defer func(begin time.Time) {
		l.log.WithFields(
			logrus.Fields{
				"duration": time.Since(begin).String(),
				"service":  "echo",
				"message":  message,
				"senderId": senderID,
			},
		).Info("SERVICE")
	}(time.Now())
	return l.s.echo(senderID, message)
}

func (l logger) sendSuggestion(id string, msg string) {
	defer func(begin time.Time) {
		l.log.WithFields(
			logrus.Fields{
				"duration":  time.Since(begin),
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
				"duration":    time.Since(begin),
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
			"duration":    time.Since(begin),
			"service":     "addToWatchlist",
			"sender_id":   senderID,
			"company_url": companyURL,
		}).Info("SERVICE")

	}(time.Now())
	l.s.addToWatchlist(senderID, companyURL)
}

func (l logger) sendWishList(senderID string) error {
	err := l.s.sendWishList(senderID)
	defer func(begin time.Time) {
		if err != nil {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).String(),
				"service":   "sendWishList",
				"sender_id": senderID,
				"error":     err.Error(),
			}).Warn("SERVICE")
		} else {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).String(),
				"service":   "sendWishList",
				"sender_id": senderID,
			}).Info("SERVICE")
		}
	}(time.Now())
	return nil
}

func (l logger) viewActiveStocks(senderID string) error {
	err := l.s.viewActiveStocks(senderID)
	defer func(begin time.Time) {
		if err != nil {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).String(),
				"service":   "viewActiveStocks",
				"sender_id": senderID,
				"error":     err.Error(),
			}).Warn("SERVICE")
		} else {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).String(),
				"service":   "viewActiveStocks",
				"sender_id": senderID,
			}).Info("SERVICE")
		}
	}(time.Now())
	return nil
}

func (l logger) editWatchList(senderID string) error {
	err := l.s.editWatchList(senderID)
	defer func(begin time.Time) {
		if err != nil {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).String(),
				"service":   "editWatchList",
				"sender_id": senderID,
				"error":     err.Error(),
			}).Warn("SERVICE")
		} else {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).String(),
				"service":   "editWatchList",
				"sender_id": senderID,
			}).Info("SERVICE")
		}
	}(time.Now())
	return nil
}
func (l logger) deleteWatchlist(senderID string, stockID string) error {
	err := l.s.deleteWatchlist(senderID, stockID)
	defer func(begin time.Time) {
		if err != nil {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).String(),
				"service":   "deleteWatchList",
				"sender_id": senderID,
				"stockID":   stockID,
				"error":     err.Error(),
			}).Warn("SERVICE")
		} else {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).String(),
				"service":   "editWatchList",
				"sender_id": senderID,
				"stock_id":  stockID,
			}).Info("SERVICE")
		}
	}(time.Now())
	return nil
}
func (l logger) sendAnnualReport(senderID string, companyURL string) error {
	defer func(begin time.Time) {
		l.log.WithFields(logrus.Fields{
			"duration":  time.Since(begin).String(),
			"service":   "sendAnnualReport",
			"sender_id": senderID,
			"stock_id":  companyURL,
		}).Info("SERVICE")
	}(time.Now())
	return l.s.sendAnnualReport(senderID, companyURL)
}

func (l logger) sendCashFlow(senderID string, companyURL string) error {
	err := l.s.sendCashFlow(senderID, companyURL)
	defer func(begin time.Time) {
		if err != nil {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).Seconds(),
				"service":   "sendCashFlow",
				"sender_id": senderID,
				"stock_id":  companyURL,
				"error":     err.Error(),
			}).Error("SERVICE")
		} else {
			l.log.WithFields(logrus.Fields{
				"duration":  time.Since(begin).Seconds(),
				"service":   "sendCashFlow",
				"sender_id": senderID,
				"stock_id":  companyURL,
			}).Info("SERVICE")
		}

	}(time.Now())

	return err
}
