// Package test drives the built hc binary and scratch image against real
// containers (testcontainers-go). Black-box: no package under test is imported.
package test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/Moq77111113/hc/test/support"
)

// TestMain builds the binary and image once. No Docker means nothing to test: exit clean.
func TestMain(m *testing.M) {
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("docker not found on PATH: skipping black-box integration tests")
		os.Exit(0)
	}

	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		log.Fatalf("set TESTCONTAINERS_RYUK_DISABLED: %v", err)
	}

	if err := support.Build(); err != nil {
		log.Fatalf("build hc artifacts: %v", err)
	}

	os.Exit(m.Run())
}
