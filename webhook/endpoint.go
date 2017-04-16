package webhook

import (
	"context"
	"fmt"

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
		} `json:"messaging"`
	} `json:"entry"`
}

// MakeEchoEndpoint creates endpoint for webhook
func makeEchoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		r := req.(request)

		result := svc.echo(r)
		return result, nil
	}
}

// EchoRequestDecoder is used to decode the request
func echoRequestDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	fmt.Print(req)

	return req, nil
}

func echoResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Write([]byte(response.(string)))
	return nil
}
