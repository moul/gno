package sessions

import (
	"testing"
)

func TestRender(t *testing.T) {
	// Test default render
	result := Render("")
	if !contains(result, "Account Sessions Demo") {
		t.Error("Expected demo information in default render")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		len(s) > len(substr) && contains(s[1:], substr)
}