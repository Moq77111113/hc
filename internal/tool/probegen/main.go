package main

import (
	"log"
)

func main() {
	if err := writeSlimTests("internal/probe"); err != nil {
		log.Fatalf("slim tests: %v", err)
	}
}
