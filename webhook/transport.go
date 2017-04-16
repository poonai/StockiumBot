package webhook

import httptransport "github.com/go-kit/kit/transport/http"
import "github.com/gorilla/mux"
import "net/http"

// MakeHandler sd
func MakeHandler(svc Service) http.Handler {
	echoHandler := httptransport.NewServer(makeEchoEndpoint(svc), echoRequestDecoder, echoResponseEncoder)
	h := mux.NewRouter()
	h.Handle("/", echoHandler).Methods("POST")
	return h
}
