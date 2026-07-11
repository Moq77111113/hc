package test

import (
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Moq77111113/hc/test/support"
)

const (
	amqpImage = "rabbitmq:4-alpine"
	amqpPort  = "5672"
)

func amqpRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        amqpImage,
		ExposedPorts: []string{amqpPort + "/tcp"},
		WaitingFor:   wait.ForLog("Server startup complete").WithStartupTimeout(2 * time.Minute),
	}
}

func TestAMQPProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, amqpRequest()), amqpPort)

	support.AssertOpenClosed(t, "amqp", addr)
}
