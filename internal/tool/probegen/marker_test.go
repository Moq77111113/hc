package main

import "testing"

func TestReplaceBlock(t *testing.T) {
	in := "a\nBEGIN\nold\nEND\nz"
	got, err := replaceBlock(in, "BEGIN", "END", "new")
	if err != nil {
		t.Fatal(err)
	}
	want := "a\nBEGIN\nnew\nEND\nz"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestReplaceBlockMissingMarker(t *testing.T) {
	if _, err := replaceBlock("no markers", "BEGIN", "END", "x"); err == nil {
		t.Fatal("expected error when markers absent")
	}
}
