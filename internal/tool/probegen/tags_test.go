package main

import (
	"testing"

	"github.com/Moq77111113/hc/internal/probe"
)

func TestSlimTags(t *testing.T) {
	got := slimTags([]probe.Scheme{probe.TCP, probe.Postgres})
	want := "hc_slim hc_tcp hc_postgres"
	if got != want {
		t.Fatalf("slimTags = %q, want %q", got, want)
	}
}
