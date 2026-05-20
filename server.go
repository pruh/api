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
	"github.com/urfave/negroni"
)

func main() {
	flag.Parse()
	err := flag.Lookup("logtostderr").Value.Set("true")
	if err != nil {
		glog.Warningf("Cannot set a flag. %s", err)
	}

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

	glog.Infof("listening on :%s", *config.Port)
	glog.Fatal(http.ListenAndServe(":"+*config.Port, router))
}
