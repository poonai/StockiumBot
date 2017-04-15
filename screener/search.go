package screener

import (
	"encoding/json"
	"net/url"

	"github.com/franela/goreq"
)

// StockSuggestion will be having the suggestion of stocks
type StockSuggestion struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

// SearchAPI for searching stocks from screener.in
const (
	SearchAPI = "https://www.screener.in/api/company/search/"
)

// SearchStock will scrape stocks from screener.in
func SearchStock(q string) ([]StockSuggestion, error) {
	query := url.Values{}
	query.Add("q", q)
	res, err := goreq.Request{
		Uri:         SearchAPI,
		QueryString: query,
		Accept:      "application/json",
	}.Do()
	if err != nil {
		return nil, err
	}
	var suggestion []StockSuggestion
	err = json.NewDecoder(res.Body).Decode(&suggestion)
	if err != nil {
		return nil, err
	}
	return suggestion, nil
}

// StockDeatail will be having the detail of stocks
type StockDeatail struct {
	Prime     string `json:"prime"`
	NumberSet struct {
		Balancesheet []struct {
			Num0 string                 `json:"0"`
			Num1 map[string]interface{} `json:"1"`
		} `json:"balancesheet"`
		Annual []struct {
			Num0 string                 `json:"0"`
			Num1 map[string]interface{} `json:"1"`
		} `json:"annual"`
		Cashflow []struct {
			Num0 string                 `json:"0"`
			Num1 map[string]interface{} `json:"1"`
		} `json:"cashflow"`
		Quarters []struct {
			Num0 string                 `json:"0"`
			Num1 map[string]interface{} `json:"1"`
		} `json:"quarters"`
	} `json:"number_set"`
	BseCode          string        `json:"bse_code"`
	ShortName        string        `json:"short_name"`
	NseCode          string        `json:"nse_code"`
	CompanyratingSet []interface{} `json:"companyrating_set"`
	AnnualreportSet  []struct {
		Source     string `json:"source"`
		ReportDate int    `json:"report_date"`
		Link       string `json:"link"`
	} `json:"annualreport_set"`
	AnnouncementSet []struct {
		AnnDate      string `json:"ann_date"`
		Announcement string `json:"announcement"`
		Link         string `json:"link"`
	} `json:"announcement_set"`
	WarehouseSet struct {
		HighPrice                    float64     `json:"high_price"`
		LowPrice                     float64     `json:"low_price"`
		SalesGrowth                  float64     `json:"sales_growth"`
		CurrentPrice                 float64     `json:"current_price"`
		DividendYield                float64     `json:"dividend_yield"`
		FaceValue                    float64     `json:"face_value"`
		ID                           int         `json:"id"`
		SalesGrowth3Years            float64     `json:"sales_growth_3years"`
		ProfitGrowth5Years           interface{} `json:"profit_growth_5years"`
		AverageReturnOnEquity3Years  float64     `json:"average_return_on_equity_3years"`
		BookValue                    float64     `json:"book_value"`
		Status                       string      `json:"status"`
		PairURL                      interface{} `json:"pair_url"`
		SalesGrowth10Years           float64     `json:"sales_growth_10years"`
		AverageReturnOnEquity10Years float64     `json:"average_return_on_equity_10years"`
		ProfitGrowth                 float64     `json:"profit_growth"`
		MarketCapitalization         float64     `json:"market_capitalization"`
		ProfitGrowth10Years          interface{} `json:"profit_growth_10years"`
		PriceToEarning               interface{} `json:"price_to_earning"`
		Industry                     string      `json:"industry"`
		Analysis                     struct {
			Remarks []interface{} `json:"remarks"`
			Bad     []string      `json:"bad"`
			Good    []string      `json:"good"`
		} `json:"analysis"`
		ResultType                  string      `json:"result_type"`
		ProfitGrowth3Years          interface{} `json:"profit_growth_3years"`
		SalesGrowth5Years           float64     `json:"sales_growth_5years"`
		ReturnOnEquity              float64     `json:"return_on_equity"`
		AverageReturnOnEquity5Years float64     `json:"average_return_on_equity_5years"`
	} `json:"warehouse_set"`
	ID   int    `json:"id"`
	Name string `json:"name"`
}
