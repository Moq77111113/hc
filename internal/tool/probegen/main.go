package main

import (
	"log"

	"github.com/Moq77111113/hc/internal/probe"
)

func main() {
	if len(probe.Catalog) == 0 {
		log.Fatal("empty catalog")
	}

	if err := writeGoreleaser(".goreleaser.yaml"); err != nil {
		log.Fatalf("goreleaser: %v", err)
	}

	if err := writeCIMatrix(".github/workflows/ci.yml"); err != nil {
		log.Fatalf("ci matrix: %v", err)
	}
}
