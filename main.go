package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"github.com/go-zoo/bone"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sch00lb0y/StockiumBot/webhook"
)

func main() {
	log := logrus.New()
	var ws webhook.Service
	ws = webhook.NewService()
	ws = webhook.NewLogger(log, ws)
	wsHandler := webhook.MakeHandler(ws)
	bone := bone.New()
	/*auth := mux.NewRouter()
	auth.HandleFunc("/", func(w http.ResponseWriter, arg2 *http.Request) {
		w.Write([]byte(arg2.URL.Query().Get("hub.challenge")))
	})
	bone.SubRoute("/webhook", auth)*/
	bone.SubRoute("/webhook", wsHandler)
	http.ListenAndServe(":80", bone)
}
