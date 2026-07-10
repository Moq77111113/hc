// Package test drives the built hc binary and scratch image against real
// containers (testcontainers-go). Black-box: no package under test is imported.
package test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var hcBinary string // host binary, built once in TestMain

const hcImage = "hc-it:latest"

// TestMain builds the binary and image once. No Docker means nothing to test: exit clean.
func TestMain(m *testing.M) {
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("docker not found on PATH: skipping black-box integration tests")
		os.Exit(0)
	}

	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		log.Fatalf("set TESTCONTAINERS_RYUK_DISABLED: %v", err)
	}

	bin, err := buildHostBinary()
	if err != nil {
		log.Fatalf("build host binary: %v", err)
	}
	hcBinary = bin

	if err := buildScratchImage(); err != nil {
		log.Fatalf("build scratch image: %v", err)
	}

	os.Exit(m.Run())
}

func buildHostBinary() (string, error) {
	dir, err := os.MkdirTemp("", "hc-it-bin-*")
	if err != nil {
		return "", err
	}

	bin := filepath.Join(dir, "hc")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build: %w\n%s", err, out)
	}
	return bin, nil
}

func buildScratchImage() error {
	cmd := exec.Command("docker", "build", "-t", hcImage, "..")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker build: %w\n%s", err, out)
	}
	return nil
}
