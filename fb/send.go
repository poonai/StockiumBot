package fb

import (
	"net/url"

	"github.com/franela/goreq"
)

// Send is used for sending the message to the appropriate user id
func Send(id string, msg string) error {
	type message struct {
		Recipient map[string]interface{} `json:"recipient"`
		Message   map[string]interface{} `json:"message"`
	}
	query := url.Values{}
	query.Add("access_token", "EAAUT2AD6M4YBANRGGFlNQ3lfqtYslHjq3vy2O56ZA18WC3B2gCid2BKYlRXWlD7EoZBFvSMejKvapCQpYlFOCeNb88eJ0NftaD1r6ypnuMcuBKPW0o35kmMc99LougKoZBoUpmoIOMnqZAEs0s0fZCOJNu3F2zJTVrite5BjKKAZDZD")
	_, err := goreq.Request{
		Uri:         "https://graph.facebook.com/v2.6/me/messages",
		QueryString: query,
		Method:      "POST",
		Accept:      "application/json",
		ContentType: "application/json",
		Body: message{
			Recipient: map[string]interface{}{"id": id},
			Message:   map[string]interface{}{"text": msg},
		},
	}.Do()

	if err != nil {
		return err
	}
	return nil
}

// StockSuggestion it'll hold the data of suggesting data
type StockSuggestion struct {
	name string
	url  string
}

func SendStockSuggestion(id string, suggestion []StockSuggestion) {

}
