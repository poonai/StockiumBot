package main

import (
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/Sirupsen/logrus"
	"github.com/go-zoo/bone"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sch00lb0y/StockiumBot/repo/mongo"
	"github.com/sch00lb0y/StockiumBot/webhook"
)

func main() {
	//fmt.Print(strconv.FormatFloat(23.23, 'f', 2, 64))
	session, err := mgo.Dial(os.Getenv("DB_URL"))
	if err != nil {
		panic(err.Error())
	}

	db := session.DB("stockiumbot")
	fbCollection := db.C("fb")
	log := logrus.New()
	var ws webhook.Service
	var wrepo webhook.Repo
	wrepo = mongo.NewRepo(fbCollection)
	ws = webhook.NewService(wrepo)
	ws = webhook.NewLogger(log, ws)
	wsHandler := webhook.MakeHandler(ws)
	bone := bone.New()
	// auth := mux.NewRouter()
	// auth.HandleFunc("/", func(w http.ResponseWriter, arg2 *http.Request) {
	// 	w.Write([]byte(arg2.URL.Query().Get("hub.challenge")))
	// })
	//bone.SubRoute("/webhook", auth)
	bone.SubRoute("/webhook", wsHandler)
	http.ListenAndServe(":80", bone)
}
