package test_working

import "testing"

func TestSimple(t *testing.T) {
	t.Log("This is a simple working test!")
	if 1+1 != 2 {
		t.Fatal("Math is broken")
	}
}