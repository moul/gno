package testutils

import (
	"std"
	"testing"
)

func TestTestAddress(t *testing.T) {
	got := TestAddress("alice")
	expected := std.Address("g1testalicealicealicealicealice")

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}

	if !got.IsValid() {
		t.Errorf("%q should be valid bech32", got)
	}
}
