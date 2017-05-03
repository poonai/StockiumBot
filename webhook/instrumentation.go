package webhook

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentation struct {
	s              Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func NewInstrumentation(counter metrics.Counter, latency metrics.Histogram, ws Service) Service {
	return instrumentation{
		ws,
		counter,
		latency,
	}
}

func (i instrumentation) echo(senderID string, message string) string {
	defer func(begin time.Time) {
		lvs := []string{"method", "echo"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())

	}(time.Now())
	return i.s.echo(senderID, message)
}

func (i instrumentation) sendSuggestion(id string, msg string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "sendSuggestion"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())

	}(time.Now())
	i.s.sendSuggestion(id, msg)
}
func (i instrumentation) sendFinancialData(id string, companyID string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "sendFinancialData"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())

	}(time.Now())
	i.s.sendFinancialData(id, companyID)
}

func (i instrumentation) addToWatchlist(senderID string, companyURL string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "addToWatchlist"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	i.s.addToWatchlist(senderID, companyURL)
}

func (i instrumentation) sendWishList(senderID string) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "sendWishList"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.s.sendWishList(senderID)
}

func (i instrumentation) viewActiveStocks(senderID string) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "sendWishList"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.s.viewActiveStocks(senderID)
}
func (i instrumentation) editWatchList(senderID string) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "editWatchList"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.s.editWatchList(senderID)
}

func (i instrumentation) deleteWatchlist(senderID string, stockID string) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "editWatchList"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.s.deleteWatchlist(senderID, stockID)
}
func (i instrumentation) sendAnnualReport(senderID string, companyURL string) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "sendAnnualReport"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.s.sendAnnualReport(senderID, companyURL)
}

func (i instrumentation) sendCashFlow(senderID string, companyURL string) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "sendCashFlow"}
		i.requestCount.With(lvs...).Add(1)
		i.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.s.sendCashFlow(senderID, companyURL)
}
