package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersionPrintsVersionAndExitsHealthy(t *testing.T) {
	var out bytes.Buffer
	if code := runVersion(&out); code != exitHealthy {
		t.Errorf("exit = %d, want %d", code, exitHealthy)
	}
	if got := out.String(); !strings.HasPrefix(got, "hc ") {
		t.Errorf("output = %q, want it to start with %q", got, "hc ")
	}
}
