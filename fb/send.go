package fb

import (
	"net/url"
	"os"

	"github.com/franela/goreq"
)

// Send is used for sending the message to the appropriate user id
func Send(id string, msg string) error {
	type message struct {
		Recipient map[string]interface{} `json:"recipient"`
		Message   map[string]interface{} `json:"message"`
	}
	query := url.Values{}
	query.Add("access_token", os.Getenv("FB_PAGE_TOKEN"))
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

// SendStockSuggestion will send the post callback button to the uesr id
func SendStockSuggestion(id string, suggestion []StockSuggestion) {

}
