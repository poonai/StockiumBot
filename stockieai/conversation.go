package stockieai

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/franela/goreq"
)

func GetResponse(q string) (string, error) {
	query := url.Values{}
	query.Add("q", q)
	res, err := goreq.Request{
		Uri:         os.Getenv("STOCKIEAI"),
		Accept:      "application/json",
		QueryString: query,
	}.Do()
	if err != nil {
		return "", err
	}
	type response struct {
		Response string `json:"response"`
	}
	var result response
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Response, nil
}
