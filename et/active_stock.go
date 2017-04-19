package et

import (
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ActiveStock wll be containing the deatails of active stocks
type ActiveStock struct {
	Name   string `json:name`
	Price  string `json:price`
	Volume string `json:volume`
}

func ActiveStocks() ([]ActiveStock, error) {
	var stocks []ActiveStock
	cmd := exec.Command("phantomjs", "et/serverRender.js", os.Getenv("ACTIVE_URL"))
	out, _ := cmd.Output()
	htmlReader := strings.NewReader(string(out))
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return stocks, err
	}

	doc.Find(".dataList").Each(func(i int, s *goquery.Selection) {
		stocks = append(stocks, ActiveStock{
			s.Find("ul:nth-child(1) > li:nth-child(1) > p:nth-child(2) > a:nth-child(1)").Text(),
			s.Find("ul:nth-child(1) > li:nth-child(2) > span:nth-child(1)").Text(),
			s.Find("ul:nth-child(1) > li:nth-child(6) > span:nth-child(1)").Text(),
		})
	})
	return stocks, nil
}
