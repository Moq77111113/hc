package main

import (
	"log"

	"github.com/Moq77111113/hc/internal/probe"
)

func main() {
	if len(probe.Catalog) == 0 {
		log.Fatal("empty catalog")
	}
	// Generation steps are wired in later tasks.
}
