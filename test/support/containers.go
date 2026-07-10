// Package support holds the plumbing shared by the black-box integration tests:
// building the hc artifacts once and driving them against real containers.
package support

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

// Start launches req, registers cleanup, and fails the test on error.
func Start(t *testing.T, req testcontainers.ContainerRequest) testcontainers.Container {
	t.Helper()
	ctr, err := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	testcontainers.CleanupContainer(t, ctr)
	if err != nil {
		t.Fatalf("start %s: %v", req.Image, err)
	}
	return ctr
}

// Endpoint returns the container's "host:port" for an exposed port.
func Endpoint(t *testing.T, ctr testcontainers.Container, port string) string {
	t.Helper()
	ep, err := ctr.PortEndpoint(context.Background(), port+"/tcp", "")
	if err != nil {
		t.Fatalf("endpoint: %v", err)
	}
	return ep
}

// OnNetwork attaches req to nw under alias, for container-to-container probing.
func OnNetwork(req testcontainers.ContainerRequest, nw *testcontainers.DockerNetwork, alias string) testcontainers.ContainerRequest {
	req.Networks = []string{nw.Name}
	req.NetworkAliases = map[string][]string{nw.Name: {alias}}
	return req
}
