package main

import (
	"flag"
	"net/http"
	"net/http/httputil"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/pruh/api/config"
	apihttp "github.com/pruh/api/http"
	"github.com/pruh/api/http/middleware"
	"github.com/pruh/api/messages"
	"github.com/pruh/api/mongo"
	"github.com/pruh/api/notifications"
	"github.com/pruh/api/providers"
	"github.com/urfave/negroni"
)

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	config, err := config.NewFromEnv()
	if err != nil {
		panic(err)
	}

	apiV1Path := "/api/v1"

	router := mux.NewRouter().StrictSlash(false)
	apiV1Router := mux.NewRouter().PathPrefix(apiV1Path).Subrouter()
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				glog.Infoln(err)
			}
			glog.Infoln(string(requestDump))

			next(w, r)
		}),
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			middleware.AuthMiddleware(w, r, next, config)
		}),
		negroni.Wrap(apiV1Router),
	)
	router.PathPrefix(apiV1Path).Handler(n)

	// messages controller
	tc := &messages.Controller{
		Config:     config,
		HTTPClient: apihttp.NewHTTPClient(),
	}
	apiV1Router.HandleFunc("/telegram/messages/send", tc.SendMessage).Methods(http.MethodPost)

	mongoClient := mongo.NewClient(config)

	// notifications controller
	repo := &notifications.Repository{
		Mongo: mongoClient,
	}
	notif := &notifications.Controller{
		Repository: repo,
	}
	apiV1Router.HandleFunc(notifications.GetPath, notif.GetAll).Methods(http.MethodGet)
	apiV1Router.HandleFunc(notifications.SingleGetPath, notif.Get).Methods(http.MethodGet)
	apiV1Router.HandleFunc(notifications.CreatePath, notif.Create).Methods(http.MethodPost)
	apiV1Router.HandleFunc(notifications.DeletePath, notif.Delete).Methods(http.MethodDelete)

	cleaner := notifications.Cleaner{
		Repository: repo,
	}
	cleaner.StartPeriodicCleaner()

	// providers controller
	provRepo := &providers.Repository{
		Mongo: mongoClient,
	}
	provController := providers.NewController(provRepo)
	apiV1Router.HandleFunc(providers.GetPath, provController.GetAll).Methods(http.MethodGet)
	apiV1Router.HandleFunc(providers.SingleGetPath, provController.Get).Methods(http.MethodGet)
	apiV1Router.HandleFunc(providers.CreatePath, provController.Create).Methods(http.MethodPost)
	apiV1Router.HandleFunc(providers.DeletePath, provController.Delete).Methods(http.MethodDelete)

	// n.Use(negroni.HandlerFunc(AuthMiddleware)) // global middleware

	glog.Infof("listening on :%s", *config.Port)
	glog.Fatal(http.ListenAndServe(":"+*config.Port, router))
}
