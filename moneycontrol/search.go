package moneycontrol

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/franela/goreq"
)

// StockSuggestion will be having the suggestion of stocks
type StockSuggestion struct {
	StockName string `json:"stock_name"`
	LinkSrc   string `json:"link_src"`
	SectorID  string `json:"sc_sector_id"`
	Sector    string `json:"sc_sector"`
}

// SearchStock will search stock from money control
func SearchStock(q string) ([]StockSuggestion, error) {
	query := url.Values{}
	query.Add("query", q)
	query.Add("type", "1")
	query.Add("format", "json")
	res, err := goreq.Request{
		Uri:         "http://www.moneycontrol.com/mccode/common/autosuggesion.php",
		QueryString: query,
		Accept:      "application/json",
	}.Do()
	if err != nil {
		return nil, err
	}
	var resJson []StockSuggestion
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return nil, err
	}
	return resJson, nil
}

// Quote will be having the price of specific stock
type Quote struct {
	Price         string `json:"price"`
	Change        string `json:"change"`
	ChangePercent string `json:"changePercent"`
	Name          string `json:"name"`
	Volume        string `json:"volume"`
}

// GetQuote will return quote by taking bse or nse code as a parameter
func GetQuote(code string) (Quote, error) {
	suggestion, err := SearchStock(code)
	if err != nil {
		return Quote{}, err
	}
	doc, err := goquery.NewDocument(suggestion[0].LinkSrc)
	if err != nil {
		return Quote{}, err
	}
	var quote Quote
	if doc.Find("#Nse_Prc_tick").Text() != "" {
		quote.Price = doc.Find("#Nse_Prc_tick").Text()
	} else {
		quote.Price = doc.Find("#Bse_Prc_tick").Text()
	}

	if doc.Find("#bse_volume > strong:nth-child(1)").Text() != "" {
		quote.Volume = doc.Find("#bse_volume > strong:nth-child(1)").Text()
	} else {
		quote.Volume = doc.Find("#nse_volume > strong:nth-child(1)").Text()
	}
	changes := strings.Split(doc.Find("#b_changetext").Text(), " ")
	//fmt.Print(changes[3])
	quote.Change = changes[1]
	quote.ChangePercent = changes[2]
	quote.ChangePercent = strings.Replace(quote.ChangePercent, "(", " ", -1)
	quote.ChangePercent = strings.Replace(quote.ChangePercent, ")", " ", -1)
	quote.Name = doc.Find("h1.b_42").Text()

	return quote, nil
}
