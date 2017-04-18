package webhook

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/sch00lb0y/StockiumBot/fb"
	"github.com/sch00lb0y/StockiumBot/screener"
	"gopkg.in/mgo.v2/bson"
)

// Service haing interface of service
type Service interface {
	echo(req request) string
	sendSuggestion(id string, msg string)
	sendFinancialData(id string, companyID string)
	addToWatchlist(senderID string, companyURL string)
}

type service struct {
	repo Repo
}

// NewService sd
func NewService(repo Repo) service {
	return service{repo}
}

func (s service) echo(req request) string {
	//fmt.Print(msg)
	//go sendSuggestion(id, msg)
	//go fb.Send(id, msg)

	return "sd"
}

func (s service) sendSuggestion(id string, msg string) {
	text := "Do you mean?"
	stocks, _ := screener.SearchStock(msg)
	var reply []fb.QuickReplie
	for i := range stocks {
		reply = append(reply,
			fb.QuickReplie{
				ContentType: "text",
				Title:       stocks[i].Name,
				Payload:     "FINANCIALDATA:" + stocks[i].Url,
			})
	}
	if len(reply) == 0 {
		text = "Sorry, I didn't understand 😰 . Could you type stock name which listed on NSE and BSE"
	}
	response := fb.Message{
		Recipient: map[string]interface{}{"id": id},
		Message:   map[string]interface{}{"text": text, "quick_replies": reply},
	}
	fb.SendStockSuggestion(response)
}

func (s service) sendFinancialData(id string, companyID string) {
	data, err := screener.GetFinancialData(companyID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"MESSAGE": err.Error(),
		}).Warn("Financial Data Error")
	}
	text := `
	        Name:    ` + data.Name + `
					HighPrice:` + strconv.FormatFloat(data.WarehouseSet.HighPrice, 'f', 2, 64) + `
					LowPrice:  ` + strconv.FormatFloat(data.WarehouseSet.LowPrice, 'f', 2, 64) + `
					CurrentPrice:` + strconv.FormatFloat(data.WarehouseSet.CurrentPrice, 'f', 2, 64) + `
					Dividend Yeild:` + strconv.FormatFloat(data.WarehouseSet.DividendYield, 'f', 2, 64) + `
					Face Value: ` + strconv.FormatFloat(data.WarehouseSet.FaceValue, 'f', 2, 64) + `
					Book Value: ` + strconv.FormatFloat(data.WarehouseSet.BookValue, 'f', 2, 64) + `
					Industry: ` + data.WarehouseSet.Industry + `
					Market Captital: ` + strconv.FormatFloat(data.WarehouseSet.MarketCapitalization, 'f', 2, 64) + `
		`
	var reply []fb.QuickReplie
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Add to Watchlist",
		Payload:     "ADDWATCHLIST:" + companyID,
	})
	response := fb.Message{
		Recipient: map[string]interface{}{"id": id},
		Message:   map[string]interface{}{"text": text, "quick_replies": reply},
	}
	fb.SendStockSuggestion(response)
}

func (s service) addToWatchlist(senderID string, companyURL string) {
	count, err := s.repo.Count(senderID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("MONGO COUNT ERROR")
		errorSend(senderID, `
			Currently we not able to accept your request.
			Our engineers working hard to solve isue.Thank You
			`)
	} else {

		if count == 0 {

			err = s.repo.Insert(Webhook{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				SenderID:  senderID,
				Portfolio: []string{companyURL},
			})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Warn("MONGO INSERT ERROR")
				errorSend(senderID, `
					Currently we not able to accept your request.
					Our engineers working hard to solve isue.Thank You
					`)
			} else {
				fb.Send(senderID, `
					Watchlist updated sucessfully`)
			}

		} else {
			err = s.repo.Update(senderID, bson.M{"$addToSet": bson.M{"portfolio": companyURL}})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Warn("MONGO UPDATE ERROR")
				errorSend(senderID, `
					Currently we not able to accept your request.
					Our engineers working hard to solve isue.Thank You
					`)
			} else {
				fb.Send(senderID, `
					Watchlist updated sucessfully`)
			}
		}

	}

}
