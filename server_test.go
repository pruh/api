package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pruh/api/config"
)

type trackingHTTPClient struct {
	called bool
}

func (c *trackingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.called = true
	return nil, errors.New("unexpected outbound request in test")
}

type blockingServer struct {
	shutdownCalled bool
	stop           chan struct{}
}

func (s *blockingServer) ListenAndServe() error {
	<-s.stop
	return http.ErrServerClosed
}

func (s *blockingServer) Shutdown(ctx context.Context) error {
	s.shutdownCalled = true
	close(s.stop)
	return nil
}

type errServer struct {
	err error
}

func (s *errServer) ListenAndServe() error {
	return s.err
}

func (s *errServer) Shutdown(ctx context.Context) error {
	return nil
}

func TestNewRouterHealthz(t *testing.T) {
	cfg := mustConfig(t, nil)
	router := newRouter(cfg, &trackingHTTPClient{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if body := w.Body.String(); body != "ok\n" {
		t.Fatalf("expected body %q, got %q", "ok\n", body)
	}
}

func TestNewRouterMessageUnauthorizedWithConfiguredCreds(t *testing.T) {
	creds := `{"admin":"password"}`
	cfg := mustConfig(t, &creds)
	router := newRouter(cfg, &trackingHTTPClient{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/telegram/messages/send",
		strings.NewReader(`{"message":"hello","chat_id":1}`))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestNewRouterMessageBadRequestWithAuth(t *testing.T) {
	creds := `{"admin":"password"}`
	cfg := mustConfig(t, &creds)
	client := &trackingHTTPClient{}
	router := newRouter(cfg, client)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/telegram/messages/send",
		strings.NewReader(`{"message":"hello","chat_id":1`))
	req.SetBasicAuth("admin", "password")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	if client.called {
		t.Fatal("did not expect outbound telegram request for malformed input")
	}
}

func TestServeUntilDoneShutsDownOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	srv := &blockingServer{stop: make(chan struct{})}
	err := serveUntilDone(ctx, srv, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !srv.shutdownCalled {
		t.Fatal("expected shutdown to be called")
	}
}

func TestServeUntilDoneReturnsListenError(t *testing.T) {
	expected := errors.New("listen failure")
	srv := &errServer{err: expected}

	err := serveUntilDone(context.Background(), srv, time.Second)
	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

func mustConfig(t *testing.T, creds *string) *config.Configuration {
	t.Helper()
	port := "8080"
	token := "test-token"

	cfg, err := config.NewFromParams(&port, &token, nil, creds)
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}
	return cfg
}
