package techpisa

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/franela/goreq"
)

// Search funs serach ticker code from techpisa
func Search(Name string) (string, error) {
	uri := os.Getenv("TECHPISASEARCH") + "?q=" + url.QueryEscape(Name)
	res, err := goreq.Request{
		Uri:    uri,
		Accept: "application/json",
	}.Do()
	if err != nil {
		return "", err
	}
	var data [][]string
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("Sorry Ticker Not Available")
	}
	return data[0][0], nil
}

func TechnicalScan(ticker string) (map[string]string, error) {
	data := make(map[string]string)
	doc, err := goquery.NewDocument(os.Getenv("TECHPISATECHNICAL") + ticker + "/")
	if err != nil {
		return nil, err
	}
	doc.Find(`div.table-responsive:nth-child(4) > table:nth-child(1) > tbody:nth-child(2) > tr`).Each(func(i int, s *goquery.Selection) {
		var key, val string
		s.Find("td").Each(func(arg1 int, sel *goquery.Selection) {
			if arg1 == 0 {
				key = sel.Text()
			}
			if arg1 == 1 {
				val = sel.Text()
			}
		})
		key = strings.Replace(key, "(?)", "", 1)
		data[key] = val
	})

	return data, nil
}
