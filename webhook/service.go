package webhook

import (
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ryanuber/columnize"
	"github.com/sch00lb0y/StockiumBot/et"
	"github.com/sch00lb0y/StockiumBot/fb"
	"github.com/sch00lb0y/StockiumBot/moneycontrol"
	"github.com/sch00lb0y/StockiumBot/screener"
	"gopkg.in/mgo.v2/bson"
)

// Service haing interface of service
type Service interface {
	echo(senderID string, message string) string
	sendSuggestion(id string, msg string)
	sendFinancialData(id string, companyID string)
	addToWatchlist(senderID string, companyURL string)
	sendWishList(senderID string) error
	viewActiveStocks(senderID string) error
}

type service struct {
	repo Repo
}

// NewService sd
func NewService(repo Repo) service {
	return service{repo}
}

func (s service) echo(senderID string, message string) string {

	go fb.Send(senderID, message)

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
	text := []string{`Name          |` + data.Name,
		`HighPrice     |` + strconv.FormatFloat(data.WarehouseSet.HighPrice, 'f', 2, 64),
		`
					 LowPrice      |` + strconv.FormatFloat(data.WarehouseSet.LowPrice, 'f', 2, 64), `
					 CurrentPrice  |` + strconv.FormatFloat(data.WarehouseSet.CurrentPrice, 'f', 2, 64), `
					 Dividend Yeild|` + strconv.FormatFloat(data.WarehouseSet.DividendYield, 'f', 2, 64), `
					 Face Value    |` + strconv.FormatFloat(data.WarehouseSet.FaceValue, 'f', 2, 64), `
					 Book Value    |` + strconv.FormatFloat(data.WarehouseSet.BookValue, 'f', 2, 64), `
					 Industry      |` + data.WarehouseSet.Industry, `
					 Market Cap    |` + strconv.FormatFloat(data.WarehouseSet.MarketCapitalization, 'f', 2, 64), `
					 `}
	var reply []fb.QuickReplie
	var code string
	if data.BseCode != "" {
		code = data.BseCode
	} else {
		code = data.NseCode
	}
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Add to Watchlist",
		Payload:     "ADDWATCHLIST:" + code,
	})
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "No",
		Payload:     "COMMANDNO:WATCHLIST",
	})
	response := fb.Message{
		Recipient: map[string]interface{}{"id": id},
		Message:   map[string]interface{}{"text": columnize.SimpleFormat(text), "quick_replies": reply},
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

func (s service) sendWishList(senderID string) error {
	var wb Webhook
	err := s.repo.Select(senderID, &wb)
	if err != nil {

		return err
	}
	var quote []moneycontrol.Quote
	var wg sync.WaitGroup
	var errs error
	for x := range wb.Portfolio {
		wg.Add(1)
		go func(index int) {
			q, err := moneycontrol.GetQuote(wb.Portfolio[index])
			if err != nil {
				errs = err
				wg.Done()
			} else {
				quote = append(quote, q)
				wg.Done()
			}
		}(x)
	}
	wg.Wait()
	if errs != nil {
		return errs
	}
	flag := 0
	output := ""
	text := []string{}
	length := len(quote)
	for x := range quote {
		flag += 1
		text = append(text, "NAME|"+quote[x].Name)
		text = append(text, "PRICE|"+quote[x].Price)
		text = append(text, "Change %|"+quote[x].ChangePercent)
		text = append(text, "Change |"+quote[x].Change)

		output += columnize.SimpleFormat(text) + "\n ------------\n"
		text = []string{}
		if (flag-length) == 0 || flag == 5 {
			length = length - flag
			response := fb.Message{
				Recipient: map[string]interface{}{"id": senderID},
				Message:   map[string]interface{}{"text": output},
			}
			go fb.SendStockSuggestion(response)
			flag = 0
			output = ""
		}
	}
	return nil
}

func (s service) viewActiveStocks(senderID string) error {
	stocks, err := et.ActiveStocks()
	if err != nil {
		return err
	}
	var text []string
	flag := 0
	output := ""
	for index := 0; index < len(stocks); index++ {
		flag += 1
		text = append(text, "NAME|"+stocks[index].Name)
		text = append(text, "PRICE|"+stocks[index].Price)
		text = append(text, "VOLUME|"+stocks[index].Volume)
		output += columnize.SimpleFormat(text) + "\n ------------\n"
		text = []string{}
		if flag == 5 {
			response := fb.Message{
				Recipient: map[string]interface{}{"id": senderID},
				Message:   map[string]interface{}{"text": output},
			}
			go fb.SendStockSuggestion(response)
			flag = 0
			output = ""
		}
	}

	return nil
}
