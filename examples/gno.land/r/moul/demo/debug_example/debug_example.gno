// Package debug_example demonstrates the use of p/moul/debug utilities
package debug_example

import (
	"strings"
	"time"

	"gno.land/p/demo/ufmt"
	"gno.land/p/moul/debug"
	"gno.land/p/moul/md"
)

// State for demonstration
var (
	counter      int
	lastAccessed time.Time
	accessLog    []string
)

func init() {
	counter = 0
	lastAccessed = time.Now()
	accessLog = []string{"System initialized"}
}

func Render(path string) string {
	// Update state
	counter++
	lastAccessed = time.Now()
	accessLog = append(accessLog, ufmt.Sprintf("Render called with path: %s", path))
	
	// Keep access log reasonable size
	if len(accessLog) > 10 {
		accessLog = accessLog[len(accessLog)-10:]
	}
	
	out := md.H1("Debug Utilities Example") + "\n\n"
	out += md.P("This realm demonstrates p/moul/debug utilities. Add ?debug=1 to the URL to see debug information.") + "\n\n"
	
	// Main content
	out += md.H2("Main Content") + "\n"
	out += md.P(ufmt.Sprintf("This page has been accessed %d times.", counter)) + "\n"
	out += md.P(ufmt.Sprintf("Current path: %s", path)) + "\n\n"
	
	// Debug information (only shown when ?debug=1 is in query)
	out += debug.Log("State Information", ufmt.Sprintf("Counter: %d, Last Access: %s", counter, lastAccessed.Format(time.RFC3339)))
	out += debug.Meta("page-views", ufmt.Sprintf("%d", counter))
	out += debug.Meta("last-accessed", lastAccessed.Format(time.RFC3339))
	
	// Show access log in debug mode
	if len(accessLog) > 0 {
		logContent := strings.Join(accessLog, "\n")
		out += debug.Log("Access Log (last 10 entries)", logContent)
	}
	
	// Debug dump of various data types
	debugData := map[string]interface{}{
		"counter":      counter,
		"path":         path,
		"logSize":      len(accessLog),
		"lastAccessed": lastAccessed.String(),
	}
	out += debug.Dump("Debug Data", debugData)
	
	out += md.H2("Instructions") + "\n"
	out += md.BulletItem("Add " + md.Code("?debug=1") + " to the URL to see debug information") + "\n"
	out += md.BulletItem("Debug logs are only visible when debug mode is enabled") + "\n"
	out += md.BulletItem("Use debug utilities to troubleshoot realm behavior") + "\n"
	
	return out
}

// Additional function that modifies state
func IncrementCounter() string {
	oldValue := counter
	counter += 10
	
	msg := ufmt.Sprintf("Counter incremented from %d to %d", oldValue, counter)
	accessLog = append(accessLog, msg)
	
	return msg
}

// Function to demonstrate error conditions
func TriggerError(shouldError bool) string {
	if shouldError {
		// In debug mode, this would show detailed error information
		errorMsg := "Intentional error for demonstration"
		accessLog = append(accessLog, "ERROR: " + errorMsg)
		return debug.Log("Error Details", errorMsg) + "\n\nError: " + errorMsg
	}
	
	return "No error triggered"
}