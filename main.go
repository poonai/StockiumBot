package main

import (
	"net/http"

	"github.com/sch00lb0y/StockiumBot/webhook"
)

func main() {
	ws := webhook.NewService()
	wsHandler := webhook.MakeHandler(ws)
	server := http.NewServeMux()
	server.Handle("/webhook/", wsHandler)
	http.Handle("/", server)
	http.ListenAndServe(":80", nil)
}
