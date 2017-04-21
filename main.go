package main

import (
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"

	"flag"
	"fmt"

	"crypto/tls"

	"github.com/Sirupsen/logrus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-zoo/bone"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/sch00lb0y/StockiumBot/repo/mongo"
	"github.com/sch00lb0y/StockiumBot/webhook"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// fmt.Print(strconv.FormatFloat(23.23, 'f', 2, 64))
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
	fieldKeys := []string{"method"}
	ws = webhook.NewService(wrepo)
	ws = webhook.NewLogger(log, ws)
	ws = webhook.NewInstrumentation(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "facebook_webhook",
		Name:      "request_count",
		Help:      "Number of requests received",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "facebook_webhook",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds",
		}, fieldKeys),
		ws)
	wsHandler := webhook.MakeHandler(ws)
	bone := bone.New()
	auth := mux.NewRouter()
	auth.HandleFunc("/", func(w http.ResponseWriter, arg2 *http.Request) {
		w.Write([]byte(arg2.URL.Query().Get("hub.challenge")))
	})

	var mode string
	var authrization bool
	flag.BoolVar(&authrization, "auth", false, "facebook webhook auth")
	flag.Parse()
	mode = os.Getenv("MODE")
	if authrization {
		bone.SubRoute("/webhook", auth)
	} else {
		bone.SubRoute("/webhook", wsHandler)
	}

	bone.SubRoute("/metrics", stdprometheus.Handler())
	if mode == "developememt" {
		http.ListenAndServe(":80", bone)
	} else {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("pagupu.in"), //your domain here
			Cache:      autocert.DirCache("certs"),          //folder for storing certificates
		}
		http.Handle("/", bone)
		server := &http.Server{
			Addr: ":443",
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}
		err = server.ListenAndServeTLS("/etc/letsencrypt/live/pagupu.in/fullchain.pem", "/etc/letsencrypt/live/pagupu.in/privkey.pem")
               http.ListenAndServe(":80",nil)
		if err != nil {
			fmt.Print(err.Error())

		}
	}

}
