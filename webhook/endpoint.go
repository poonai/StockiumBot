package webhook

import (
	"context"
	"fmt"
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
		for _, x := range r.Entry {
			for _, y := range x.Messaging {
				if len(y.Message.QuickReply) > 0 {
					quickReply := y.Message.QuickReply
					payload := quickReply["payload"].(string)
					sep := strings.Split(payload, ":")
					switch sep[0] {
					case "FINANCIALDATA":

						go svc.sendFinancialData(y.Sender.ID, sep[1])
						break
					case "ADDWATCHLIST":
						go svc.addToWatchlist(y.Sender.ID, sep[1])
						break
					case "COMMANDNO":
						go svc.echo(y.Sender.ID, `it's okay,still you can serach stock`)
						break
					case "REMOVEWATCHLIST":
						go svc.deleteWatchlist(y.Sender.ID, sep[1])
					}
				} else {
					if y.PostBack.Payload != "" {
						payload := y.PostBack.Payload
						sep := strings.Split(payload, ":")
						if sep[0] == "COMMAND" {
							switch sep[1] {
							case "VIEWWATCHLIST":
								svc.sendWishList(y.Sender.ID)
								break
							case "VIEWACTIVESTOCKS":
								svc.viewActiveStocks(y.Sender.ID)
								break
							case "EDITWATCHLIST":
								svc.editWatchList(y.Sender.ID)
							}
						}
					} else {
						go svc.sendSuggestion(y.Sender.ID, y.Message.Text)
					}

				}

			}
		}

		//result := svc.echo(r)
		return "sup", nil
	}
}

// EchoRequestDecoder is used to decode the request
func echoRequestDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	var req request
	fmt.Print("hitting")
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
