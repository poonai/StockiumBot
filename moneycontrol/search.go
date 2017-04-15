package moneycontrol

import (
	"encoding/json"
	"net/url"

	"github.com/franela/goreq"
)

// StockSuggestion will be having the suggestion of stocks
type StockSuggestion struct {
	StockName string `json:"stock_name"`
	LinkSrc   string `json:"link_src"`
	SectorID  string `json:"sc_sector_id"`
	Sector    string `json:"sc_sector"`
}

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
