package moneycontrol

import (
	"encoding/json"
	"fmt"
	"net/url"

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
}

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
	}
	fmt.Print("hello")
	fmt.Print(doc.Find("div #n_changetext").Text())
	return Quote{}, nil
}
