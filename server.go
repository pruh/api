package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
	"github.com/pruh/api/controllers"
	"github.com/pruh/api/middleware"
	"github.com/pruh/api/utils"
	"github.com/urfave/negroni"
)

func main() {
	config, err := utils.NewFromEnv()
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
				log.Println(err)
			}
			log.Println(string(requestDump))

			next(w, r)
		}),
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			middleware.AuthMiddleware(w, r, next, config)
		}),
		negroni.Wrap(apiV1Router),
	)
	router.PathPrefix(apiV1Path).Handler(n)

	// messages controller
	tc := &controllers.TelegramController{
		Config:     config,
		HTTPClient: utils.NewHTTPClient(),
	}
	apiV1Router.HandleFunc("/telegram/messages/send", tc.SendMessage).Methods(http.MethodPost)

	// notifications controller
	notif := &controllers.NotificationsController{}
	apiV1Router.HandleFunc("/notifications", notif.Get).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/notifications", notif.Create).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/notifications", notif.Delete).Methods(http.MethodDelete)

	// n.Use(negroni.HandlerFunc(AuthMiddleware)) // global middleware

	log.Printf("listening on :%s", *config.Port)
	log.Fatal(http.ListenAndServe(":"+*config.Port, router))
}
