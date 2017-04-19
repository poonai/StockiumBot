package webhook

import (
	"context"
	"strings"

	"net/http"

	"encoding/json"

	"github.com/go-kit/kit/endpoint"
)

type request struct {
	Entry []struct {
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Message struct {
				Text       string                 `json:"text"`
				QuickReply map[string]interface{} `json:"quick_reply"`
			} `json:"message"`
			PostBack struct {
				Payload string `json:"payload"`
			} `json:"postback"`
		} `json:"messaging"`
	} `json:"entry"`
}

// MakeEchoEndpoint creates endpoint for webhook
func makeEchoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(request)
		if len(r.Entry[0].Messaging[0].Message.QuickReply) > 0 {
			quickReply := r.Entry[0].Messaging[0].Message.QuickReply
			payload := quickReply["payload"].(string)
			sep := strings.Split(payload, ":")
			switch sep[0] {
			case "FINANCIALDATA":

				go svc.sendFinancialData(r.Entry[0].Messaging[0].Sender.ID, sep[1])
				break
			case "ADDWATCHLIST":
				go svc.addToWatchlist(r.Entry[0].Messaging[0].Sender.ID, sep[1])
				break
			case "COMMANDNO":
				go svc.echo(r.Entry[0].Messaging[0].Sender.ID, `it's okay,still you can serach stock`)
				break
			}
		} else {
			if r.Entry[0].Messaging[0].PostBack.Payload != "" {
				payload := r.Entry[0].Messaging[0].PostBack.Payload
				sep := strings.Split(payload, ":")
				if sep[0] == "COMMAND" {
					switch sep[1] {
					case "VIEWWATCHLIST":
						svc.sendWishList(r.Entry[0].Messaging[0].Sender.ID)
						break
					case "VIEWACTIVESTOCKS":
						svc.viewActiveStocks(r.Entry[0].Messaging[0].Sender.ID)
					}
				}
			} else {
				go svc.sendSuggestion(r.Entry[0].Messaging[0].Sender.ID, r.Entry[0].Messaging[0].Message.Text)
			}

		}

		//result := svc.echo(r)
		return "sup", nil
	}
}

// EchoRequestDecoder is used to decode the request
func echoRequestDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func echoResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Write([]byte(response.(string)))
	return nil
}
