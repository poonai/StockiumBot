package fb

import (
	"fmt"
	"net/url"
	"os"

	"github.com/franela/goreq"
)

type Message struct {
	Recipient map[string]interface{} `json:"recipient"`
	Message   map[string]interface{} `json:"message"`
}

// Send is used for sending the message to the appropriate user id
func Send(id string, msg string) error {

	query := url.Values{}
	query.Add("access_token", os.Getenv("FB_PAGE_TOKEN"))
	_, err := goreq.Request{
		Uri:         "https://graph.facebook.com/v2.6/me/messages",
		QueryString: query,
		Method:      "POST",
		Accept:      "application/json",
		ContentType: "application/json",
		Body: Message{
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
type QuickReplie struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload"`
}

// SendStockSuggestion will send the post callback button to the uesr id
func SendStockSuggestion(msg Message) {
	query := url.Values{}
	query.Add("access_token", os.Getenv("FB_PAGE_TOKEN"))
	r, err := goreq.Request{
		Uri:         "https://graph.facebook.com/v2.6/me/messages",
		QueryString: query,
		Method:      "POST",
		Body:        msg,
		Accept:      "application/json",
		ContentType: "application/json",
	}.Do()
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(r.Body.ToString())
}
