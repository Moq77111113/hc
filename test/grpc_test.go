package test

import (
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Moq77111113/hc/test/support"
)

const (
	grpcImage = "gcr.io/etcd-development/etcd:v3.5.17"
	grpcPort  = "2379"
)

// etcd runs a gRPC server that registers the standard grpc.health.v1.Health
// service and reports SERVING, so it exercises the real Health/Check path.
func grpcRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        grpcImage,
		ExposedPorts: []string{grpcPort + "/tcp"},
		Cmd: []string{
			"etcd",
			"--listen-client-urls", "http://0.0.0.0:2379",
			"--advertise-client-urls", "http://0.0.0.0:2379",
		},
		WaitingFor: wait.ForLog("ready to serve client requests").WithStartupTimeout(2 * time.Minute),
	}
}

func TestGRPCProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, grpcRequest()), grpcPort)

	support.AssertOpenClosed(t, "grpc", addr)
}
