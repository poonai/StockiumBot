package webhook

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ryanuber/columnize"
	"github.com/sch00lb0y/StockiumBot/et"
	"github.com/sch00lb0y/StockiumBot/fb"
	"github.com/sch00lb0y/StockiumBot/moneycontrol"
	"github.com/sch00lb0y/StockiumBot/screener"
	"github.com/sch00lb0y/StockiumBot/stockieai"
	"github.com/sch00lb0y/StockiumBot/techpisa"
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
	editWatchList(senderID string) error
	deleteWatchlist(senderID string, stockID string) error
	sendAnnualReport(senderID string, companyURL string) error
	sendCashFlow(senderID string, companyURL string) error
	sendTechnicalScan(senderID string, ticker string) error
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
	convo := map[string]bool{
		"hi":   true,
		"K":    true,
		"okay": true,
	}
	msg = strings.ToLower(msg)
	if convo[msg] {
		if msg == "k" || msg == "okay" {
			fb.Send(id, "Feel Free to ask any thing")
		} else {
			res, err := stockieai.GetResponse(msg)
			if err != nil {
				fb.Send(id, "Type some stock name to get info")
			} else {
				fb.Send(id, res)
			}
		}
	} else {
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
			res, err := stockieai.GetResponse(msg)
			if err != nil {
				text = "type some stock to get info of stock"
			} else {
				text = res
			}
		}
		response := fb.Message{
			Recipient: map[string]interface{}{"id": id},
			Message:   map[string]interface{}{"text": text, "quick_replies": reply},
		}
		fb.SendStockSuggestion(response)
	}

}

func sortKeys(maps map[string]interface{}) []string {
	keys := make([]int, len(maps))
	i := 0
	for k := range maps {
		keys[i], _ = strconv.Atoi(strings.Split(k, "-")[0])
		i++
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	final := make([]string, 0)
	for i, _ := range keys {
		final = append(final, strconv.Itoa(keys[i])+"-03-31")
	}

	return final
}

func (s service) sendFinancialData(id string, companyID string) {
	data, err := screener.GetFinancialData(companyID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"MESSAGE": err.Error(),
		}).Warn("Financial Data Error")
	}

	text := []string{
		`Name          |` + data.Name,
		`HighPrice     |` + strconv.FormatFloat(data.WarehouseSet.HighPrice, 'f', 2, 64), `
					  LowPrice      |` + strconv.FormatFloat(data.WarehouseSet.LowPrice, 'f', 2, 64), `
					  CurrentPrice  |` + strconv.FormatFloat(data.WarehouseSet.CurrentPrice, 'f', 2, 64), `
					  Dividend Yeild|` + strconv.FormatFloat(data.WarehouseSet.DividendYield, 'f', 2, 64), `
					  Face Value    |` + strconv.FormatFloat(data.WarehouseSet.FaceValue, 'f', 2, 64), `
				 	  Book Value    |` + strconv.FormatFloat(data.WarehouseSet.BookValue, 'f', 2, 64), `
					  Industry      |` + data.WarehouseSet.Industry, `
		        Market Cap    |` + strconv.FormatFloat(data.WarehouseSet.MarketCapitalization, 'f', 2, 64), `
					 Avg ROE 5 Years|` + strconv.FormatFloat(data.WarehouseSet.AverageReturnOnEquity5Years, 'f', 2, 64)}
	f, ok := data.WarehouseSet.ProfitGrowth5Years.(float64)
	if ok {
		text = append(text, "ProfitGrowth 5 Years |"+strconv.FormatFloat(f, 'f', 2, 64))

	}
	f, ok = data.WarehouseSet.ProfitGrowth3Years.(float64)
	if ok {
		text = append(text, "Profit Growth 3 Yrs"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	f, ok = data.WarehouseSet.ProfitGrowth10Years.(float64)
	if ok {
		text = append(text, "Profit Growth 3 Yrs"+strconv.FormatFloat(f, 'f', 2, 64))
	}

	var reply []fb.QuickReplie
	var code string
	if data.BseCode != "" {
		code = data.BseCode
	} else {
		code = data.NseCode
	}
	reply = financialQuoteReply(companyID, "FINANCIALDATA")
	ticker := ""
	if data.NseCode != "" {
		ticker = data.NseCode
	} else {
		ticker = "podadey"
	}
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Add to Watchlist",
		Payload:     "ADDWATCHLIST:" + code,
	})
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Technical Scan",
		Payload:     "TECHSCAN:" + ticker,
	})
	response := fb.Message{
		Recipient: map[string]interface{}{"id": id},
		Message:   map[string]interface{}{"text": columnize.SimpleFormat(text), "quick_replies": reply},
	}
	fb.SendStockSuggestion(response)
	// totalLiabilites := data.NumberSet.Balancesheet[4].([]interface{})[1].(map[string]interface{})
	// keys := sortKeys(totalLiabilites)
	// text = []string{}
	// text = append(text, "Year|Amt")
	// f, ok = totalLiabilites[keys[2]].(float64)
	// if ok {
	// 	text = append(text, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = totalLiabilites[keys[1]].(float64)
	// if ok {
	// 	text = append(text, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = totalLiabilites[keys[0]].(float64)
	// if ok {
	// 	text = append(text, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// totalAsset := data.NumberSet.Balancesheet[9].([]interface{})[1].(map[string]interface{})
	// asset := []string{}
	// asset = append(asset, "Year|Amt")
	// keys = sortKeys(totalAsset)
	// f, ok = totalAsset[keys[2]].(float64)
	// if ok {
	// 	asset = append(asset, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = totalAsset[keys[1]].(float64)
	// if ok {
	// 	asset = append(asset, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = totalAsset[keys[0]].(float64)
	// if ok {
	// 	asset = append(asset, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// response = fb.Message{
	// 	Recipient: map[string]interface{}{"id": id},
	// 	Message: map[string]interface{}{"text": "\nBalance Sheet\n-------------\nTotal Liabilities \n-------------------\n" + columnize.SimpleFormat(text) +
	// 		"\n Total Asset \n-------------\n" + columnize.SimpleFormat(asset),
	// 	}, //"quick_replies": reply},
	// }
	//	fb.SendStockSuggestion(response)
	// OperationProfit := data.NumberSet.Annual[2].([]interface{})[1].(map[string]interface{})
	// keys = sortKeys(OperationProfit)
	// op := []string{}
	// op = append(op, "Year|Amt")
	// f, ok = OperationProfit[keys[2]].(float64)
	// if ok {
	// 	op = append(op, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = OperationProfit[keys[1]].(float64)
	// if ok {
	// 	op = append(op, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = OperationProfit[keys[0]].(float64)
	// if ok {
	// 	op = append(op, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// ProfitBeforeTax := data.NumberSet.Annual[7].([]interface{})[1].(map[string]interface{})
	// keys = sortKeys(ProfitBeforeTax)
	// pbt := []string{}
	// pbt = append(pbt, "Year|Amt")
	// f, ok = ProfitBeforeTax[keys[2]].(float64)
	// if ok {
	// 	pbt = append(pbt, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = ProfitBeforeTax[keys[1]].(float64)
	// if ok {
	// 	pbt = append(pbt, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = ProfitBeforeTax[keys[0]].(float64)
	// if ok {
	// 	pbt = append(pbt, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// NetProfit := data.NumberSet.Annual[9].([]interface{})[1].(map[string]interface{})
	// keys = sortKeys(NetProfit)
	// np := []string{}
	// np = append(np, "Year|Amt")
	// f, ok = NetProfit[keys[2]].(float64)
	// if ok {
	// 	np = append(np, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = NetProfit[keys[1]].(float64)
	// if ok {
	// 	np = append(np, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = NetProfit[keys[0]].(float64)
	// if ok {
	// 	np = append(np, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// CashFlow := data.NumberSet.Cashflow[3].([]interface{})[1].(map[string]interface{})
	// keys = sortKeys(CashFlow)
	// cf := []string{}
	// fmt.Print(CashFlow)
	// cf = append(cf, "Year|Amt")
	// f, ok = CashFlow[keys[2]].(float64)
	// if ok {
	// 	cf = append(cf, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = CashFlow[keys[1]].(float64)
	// if ok {
	// 	cf = append(cf, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// f, ok = NetProfit[keys[0]].(float64)
	// if ok {
	// 	cf = append(cf, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	// }
	// response = fb.Message{
	// 	Recipient: map[string]interface{}{"id": id},
	// 	Message: map[string]interface{}{"text": "\nAnnual Report\n-------------\nOperational Profit \n-------------------\n" + columnize.SimpleFormat(op) +
	// 		"\nProfit Before Tax\n----------\n" + columnize.SimpleFormat(pbt) +
	// 		"\n Net Profit\n-----------\n" + columnize.SimpleFormat(np) +
	// 		"\n Cash FLow\n----------\n" + columnize.SimpleFormat(cf),
	// 		"quick_replies": reply},
	// }
	// fb.SendStockSuggestion(response)
}

func (s service) sendCashFlow(senderID string, companyURL string) error {
	data, err := screener.GetFinancialData(companyURL)
	if err != nil {
		return err
	}
	CashFlow := data.NumberSet.Cashflow[3].([]interface{})[1].(map[string]interface{})
	keys := sortKeys(CashFlow)
	cf := []string{}

	cf = append(cf, "Year|Amt")
	f, ok := CashFlow[keys[2]].(float64)
	if ok {
		cf = append(cf, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	f, ok = CashFlow[keys[1]].(float64)
	if ok {
		cf = append(cf, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	var reply []fb.QuickReplie
	var code string
	if data.BseCode != "" {
		code = data.BseCode
	} else {
		code = data.NseCode
	}
	reply = financialQuoteReply(companyURL, "CASHFLOW")
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Add to Watchlist",
		Payload:     "ADDWATCHLIST:" + code,
	})
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Technical Scan",
		Payload:     "TECHSCAN:" + data.Name,
	})

	response := fb.Message{
		Recipient: map[string]interface{}{"id": senderID},
		Message: map[string]interface{}{"text": "\n Cash FLow\n----------\n" + columnize.SimpleFormat(cf),
			"quick_replies": reply},
	}
	fb.SendStockSuggestion(response)
	return nil
}

func (s service) sendAnnualReport(senderID string, companyURL string) error {
	data, err := screener.GetFinancialData(companyURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"MESSAGE": err.Error(),
		}).Warn("Financial Data Error")
	}

	OperationProfit := data.NumberSet.Annual[2].([]interface{})[1].(map[string]interface{})
	keys := sortKeys(OperationProfit)
	op := []string{}
	op = append(op, "Year|Amt")
	f, ok := OperationProfit[keys[2]].(float64)
	if ok {
		op = append(op, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	f, ok = OperationProfit[keys[1]].(float64)
	if ok {
		op = append(op, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	f, ok = OperationProfit[keys[0]].(float64)
	if ok {
		op = append(op, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	ProfitBeforeTax := data.NumberSet.Annual[7].([]interface{})[1].(map[string]interface{})
	keys = sortKeys(ProfitBeforeTax)
	pbt := []string{}
	pbt = append(pbt, "Year|Amt")
	f, ok = ProfitBeforeTax[keys[2]].(float64)
	if ok {
		pbt = append(pbt, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	f, ok = ProfitBeforeTax[keys[1]].(float64)
	if ok {
		pbt = append(pbt, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	f, ok = ProfitBeforeTax[keys[0]].(float64)
	if ok {
		pbt = append(pbt, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	NetProfit := data.NumberSet.Annual[9].([]interface{})[1].(map[string]interface{})
	keys = sortKeys(NetProfit)
	np := []string{}
	np = append(np, "Year|Amt")
	f, ok = NetProfit[keys[2]].(float64)
	if ok {
		np = append(np, keys[2]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	f, ok = NetProfit[keys[1]].(float64)
	if ok {
		np = append(np, keys[1]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	f, ok = NetProfit[keys[0]].(float64)
	if ok {
		np = append(np, keys[0]+"|"+strconv.FormatFloat(f, 'f', 2, 64))
	}
	var reply []fb.QuickReplie
	var code string
	if data.BseCode != "" {
		code = data.BseCode
	} else {
		code = data.NseCode
	}
	reply = financialQuoteReply(companyURL, "ANNUALREPORT")
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Add to Watchlist",
		Payload:     "ADDWATCHLIST:" + code,
	})
	ticker := ""
	if data.NseCode != "" {
		ticker = data.NseCode
	} else {
		ticker = "podadey"
	}
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Technical Scan",
		Payload:     "TECHSCAN:" + ticker,
	})
	response := fb.Message{
		Recipient: map[string]interface{}{"id": senderID},
		Message: map[string]interface{}{"text": "\nAnnual Report\n-------------\nOperational Profit \n-------------------\n" + columnize.SimpleFormat(op) +
			"\nProfit Before Tax\n----------\n" + columnize.SimpleFormat(pbt) +
			"\n Net Profit\n-----------\n" + columnize.SimpleFormat(np),
			"quick_replies": reply},
	}
	fb.SendStockSuggestion(response)
	return nil
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
					Our engineers working hard to solve issue. Thank You
					`)
			} else {
				fb.Send(senderID, `
					Watchlist updated sucessfully`)
			}
		}

	}
}

func financialQuoteReply(ticker string, statement string) []fb.QuickReplie {
	repies := map[string]fb.QuickReplie{
		"CASHFLOW": fb.QuickReplie{
			ContentType: "text",
			Title:       "Cash Flow",
			Payload:     "CASHFLOW:" + ticker,
		},
		"ANNUALREPORT": fb.QuickReplie{
			ContentType: "text",
			Title:       "Annual Report",
			Payload:     "ANNUALREPORT:" + ticker,
		},
		"FINANCIALDATA": fb.QuickReplie{
			ContentType: "text",
			Title:       "Financial Data",
			Payload:     "FINANCIALDATA:" + ticker,
		},
		"COMMANDNO": fb.QuickReplie{
			ContentType: "text",
			Title:       "No",
			Payload:     "COMMANDNO:WATCHLIST",
		},
	}
	var reply []fb.QuickReplie
	for i := range repies {
		if i != statement {
			reply = append(reply, repies[i])
		}
	}
	return reply
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
	if len(wb.Portfolio) > 0 {
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
			text = append(text, "Change |"+quote[x].ChangePercent)
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

	} else {
		fb.Send(senderID, "You don't have stock in watchlist to see")
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
func (s service) editWatchList(senderID string) error {
	var wb Webhook
	err := s.repo.Select(senderID, &wb)
	if err != nil {
		return err
	}
	type temp struct {
		Name string
		ID   string
	}
	quote := []temp{}
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
				quote = append(quote, temp{
					Name: q.Name,
					ID:   wb.Portfolio[index],
				})
				wg.Done()
			}
		}(x)
	}
	wg.Wait()
	if errs != nil {
		return errs
	}
	var reply []fb.QuickReplie
	for i := range quote {
		reply = append(reply,
			fb.QuickReplie{
				ContentType: "text",
				Title:       quote[i].Name,
				Payload:     "REMOVEWATCHLIST:" + quote[i].ID,
			})
	}
	reply = append(reply, fb.QuickReplie{
		ContentType: "text",
		Title:       "Cancel",
		Payload:     "COMMANDNO:WATCHLIST",
	})
	response := fb.Message{
		Recipient: map[string]interface{}{"id": senderID},
		Message: map[string]interface{}{"text": "Select stock to remove",
			"quick_replies": reply,
		},
	}
	fb.SendStockSuggestion(response)
	return nil
}

func (s service) deleteWatchlist(senderID string, stockID string) error {
	err := s.repo.Update(senderID, bson.M{"$pull": bson.M{"portfolio": stockID}})
	response := fb.Message{}
	response.Recipient = map[string]interface{}{"id": senderID}
	if err != nil {
		response.Message = map[string]interface{}{"text": "something went wrong our engineers working hard to find the issue"}
		fb.SendStockSuggestion(response)
		return err
	}
	response.Message = map[string]interface{}{"text": "Updated SucessFully"}
	fb.SendStockSuggestion(response)
	return nil
}

func (s service) sendTechnicalScan(senderID string, stockName string) error {
	ticker, err := techpisa.Search(stockName)
	if err != nil {
		fb.Send(senderID, "sorry scans are not available")
		return err
	}
	result, err := techpisa.TechnicalScan(ticker)

	if err != nil {
		fb.Send(senderID, "sorry scans are not available")
		return err
	}
	var text string
	for key, val := range result {
		text += "\n--" + key + "--\n"
		if len(text+val) > 640 {
			response := fb.Message{
				Recipient: map[string]interface{}{"id": senderID},
				Message:   map[string]interface{}{"text": text},
			}
			fb.SendStockSuggestion(response)
			text = ""
		}
		text += val
	}
	if text != "" {
		response := fb.Message{
			Recipient: map[string]interface{}{"id": senderID},
			Message:   map[string]interface{}{"text": text},
		}
		fb.SendStockSuggestion(response)
	}

	return nil
}
