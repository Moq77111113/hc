package test

import (
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Moq77111113/hc/test/support"
)

const (
	mysqlImage = "mysql:8"
	mysqlPort  = "3306"
)

func mysqlRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        mysqlImage,
		ExposedPorts: []string{mysqlPort + "/tcp"},
		Env:          map[string]string{"MYSQL_ROOT_PASSWORD": "pw"},
		WaitingFor:   wait.ForExposedPort().WithStartupTimeout(2 * time.Minute),
	}
}

func TestMySQLProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, mysqlRequest()), mysqlPort)

	support.AssertOpenClosed(t, "mysql", addr)
}
