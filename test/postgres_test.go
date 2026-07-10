package test

import (
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

func TestPostgresProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, pgRequest()), "5432")

	assertOpenClosed(t, "postgres", addr)
}
