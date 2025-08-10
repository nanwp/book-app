package service

import (
	"context"
	"testing"
)

func TestProcessURL_All(t *testing.T) {
	ctx := context.TODO()
	in := "https://BYFOOD.com/food-EXPeriences?query=abc/"
	want := "https://www.byfood.com/food-experiences"
	got, err := ProcessURL(ctx, in, "all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestProcessURL_Canonical(t *testing.T) {
	ctx := context.TODO()
	in := "https://BYFOOD.com/food-EXPeriences?query=abc/"
	want := "https://BYFOOD.com/food-EXPeriences"
	got, err := ProcessURL(ctx, in, "canonical")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestProcessURL_Redirection(t *testing.T) {
	ctx := context.TODO()
	in := "https://BYFOOD.com/Food/Bar?x=1&y=2"
	want := "https://www.byfood.com/food/bar?x=1&y=2" // query stays; whole URL lowercased
	got, err := ProcessURL(ctx, in, "redirection")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestProcessURL_Validation(t *testing.T) {
	ctx := context.TODO()
	cases := []struct {
		name string
		url  string
		op   Operation
	}{
		{"no scheme", "//example.com/path", OpAll},
		{"bad scheme", "ftp://example.com/path", OpAll},
		{"no host", "https:///path", OpAll},
	}
	for _, tc := range cases {
		if _, err := ProcessURL(ctx, tc.url, tc.op.String()); err == nil {
			t.Fatalf("%s: expected error, got nil", tc.name)
		}
	}
}

func TestParseOperation(t *testing.T) {
	if _, err := ParseOperation("bad"); err == nil {
		t.Fatalf("expected error for bad operation")
	}
	if op, err := ParseOperation("CANONICAL"); err != nil || op != OpCanonical {
		t.Fatalf("parse canonical failed: %v %v", op, err)
	}
}
