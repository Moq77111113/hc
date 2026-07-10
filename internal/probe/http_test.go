//go:build !hc_slim || hc_http

package probe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHTTPProberHealthyOn2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	if err := (httpProber{}).Probe(context.Background(), u); err != nil {
		t.Fatalf("want healthy, got %v", err)
	}
}

func TestHTTPProberUnhealthyOn5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	if err := (httpProber{}).Probe(context.Background(), u); err == nil {
		t.Fatal("want error on 500, got nil")
	}
}
