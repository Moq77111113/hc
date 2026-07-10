package test

import (
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Moq77111113/hc/test/support"
)

const (
	pgImage = "postgres:18-alpine"
	pgPort  = "5432"
)

func pgRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        pgImage,
		ExposedPorts: []string{pgPort + "/tcp"},
		Env:          map[string]string{"POSTGRES_PASSWORD": "pw"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
}

func TestPostgresProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, pgRequest()), pgPort)

	support.AssertOpenClosed(t, "postgres", addr)
}
