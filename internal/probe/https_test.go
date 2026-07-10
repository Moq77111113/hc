//go:build !hc_slim || hc_https

package probe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHTTPSProberHealthyOn2xx(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	if err := (httpsProber{}).Probe(context.Background(), u); err != nil {
		t.Fatalf("want healthy (skip-verify), got %v", err)
	}
}

func TestHTTPSProberUnhealthyOn5xx(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	if err := (httpsProber{}).Probe(context.Background(), u); err == nil {
		t.Fatal("want error on 500, got nil")
	}
}
