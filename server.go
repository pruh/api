package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/pruh/api/v3/config"
	apihttp "github.com/pruh/api/v3/http"
	"github.com/pruh/api/v3/http/middleware"
	"github.com/pruh/api/v3/messages"
	"github.com/urfave/negroni"
)

type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

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

	router := newRouter(config, apihttp.NewHTTPClient())
	srv := &http.Server{
		Addr:    ":" + *config.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	glog.Infof("listening on :%s", *config.Port)
	if err := serveUntilDone(ctx, srv, 10*time.Second); err != nil {
		glog.Fatalf("server error: %v", err)
	}
}

func newRouter(config *config.Configuration, httpClient apihttp.Client) *mux.Router {
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
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	}).Methods(http.MethodGet)

	// messages controller
	tc := &messages.Controller{
		Config:     config,
		HTTPClient: httpClient,
	}
	apiV1Router.HandleFunc("/telegram/messages/send", tc.SendMessage).Methods(http.MethodPost)

	return router
}

func serveUntilDone(ctx context.Context, srv server, shutdownTimeout time.Duration) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
	}

	glog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	err := <-errCh
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
