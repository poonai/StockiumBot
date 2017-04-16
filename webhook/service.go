package webhook

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/sch00lb0y/StockiumBot/fb"
	"github.com/sch00lb0y/StockiumBot/screener"
)

// Service haing interface of service
type Service interface {
	echo(req request) string
}

type service struct {
}

// NewService sd
func NewService() service {
	return service{}
}

func (s service) echo(req request) string {
	//fmt.Print(msg)
	//go sendSuggestion(id, msg)
	//go fb.Send(id, msg)

	if len(req.Entry[0].Messaging[0].Message.QuickReply) > 0 {
		quickReply := req.Entry[0].Messaging[0].Message.QuickReply
		payload := quickReply["payload"].(string)
		sep := strings.Split(payload, ":")
		switch sep[0] {
		case "FINANCIALDATA":
			fmt.Print(sep)
			go sendFinancialData(req.Entry[0].Messaging[0].Sender.ID, sep[1])
			break
		}
	} else {
		sendSuggestion(req.Entry[0].Messaging[0].Sender.ID, req.Entry[0].Messaging[0].Message.Text)
	}
	return "sd"
}

func sendSuggestion(id string, msg string) {
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

func sendFinancialData(id string, companyID string) {
	data, err := screener.GetFinancialData(companyID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"MESSAGE": err.Error(),
		}).Warn("Financial Data Error")
	}
	text := `
	        Name:    ` + data.Name + `
					HighPrice:` + strconv.FormatFloat(data.WarehouseSet.HighPrice, 'f', 6, 64) + `
					LowPrice:  ` + strconv.FormatFloat(data.WarehouseSet.LowPrice, 'f', 6, 64) + `
					CurrentPrice:` + strconv.FormatFloat(data.WarehouseSet.CurrentPrice, 'f', 6, 64) + `
					Dividend Yeild:` + strconv.FormatFloat(data.WarehouseSet.DividendYield, 'f', 6, 64) + `
					Face Value: ` + strconv.FormatFloat(data.WarehouseSet.FaceValue, 'f', 6, 64) + `
					Book Value: ` + strconv.FormatFloat(data.WarehouseSet.BookValue, 'f', 6, 64) + `
					Industry: ` + data.WarehouseSet.Industry + `
					Market Captital: ` + strconv.FormatFloat(data.WarehouseSet.MarketCapitalization, 'f', 6, 64) + `
		`
	response := fb.Message{
		Recipient: map[string]interface{}{"id": id},
		Message:   map[string]interface{}{"text": text},
	}
	fb.SendStockSuggestion(response)
}
