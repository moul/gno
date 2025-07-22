package forms

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	// Test homepage
	homepage := Render("")
	if !strings.Contains(homepage, "Gno Forms Documentation") {
		t.Error("Homepage should contain title")
	}
	if !strings.Contains(homepage, "Available Examples") {
		t.Error("Homepage should list examples")
	}
	
	// Test simple form page
	simple := Render("simple")
	if !strings.Contains(simple, "Simple Form Example") {
		t.Error("Simple form page should have correct title")
	}
	if !strings.Contains(simple, "<gno-form") {
		t.Error("Simple form page should contain a form")
	}
	
	// Test registration form
	reg := Render("registration")
	if !strings.Contains(reg, "Registration Form Example") {
		t.Error("Registration page should have correct title")
	}
	
	// Test form submission handling
	submission := Render("test?name=John&email=john@example.com")
	if !strings.Contains(submission, "Form Submitted Successfully") {
		t.Error("Submission should show success message")
	}
}

func TestParseQueryParams(t *testing.T) {
	params := parseQueryParams("name=John&email=test@example.com&age=25")
	
	if params["name"] != "John" {
		t.Errorf("Expected name=John, got %s", params["name"])
	}
	if params["email"] != "test@example.com" {
		t.Errorf("Expected email=test@example.com, got %s", params["email"])
	}
	if params["age"] != "25" {
		t.Errorf("Expected age=25, got %s", params["age"])
	}
}

func TestFormatSubmissionData(t *testing.T) {
	fields := map[string]string{
		"username": "john",
		"password": "secret123",
		"email":    "john@example.com",
	}
	
	formatted := formatSubmissionData(fields)
	
	// Check that password is masked
	if strings.Contains(formatted, "secret123") {
		t.Error("Password should be masked in output")
	}
	if !strings.Contains(formatted, "********") {
		t.Error("Password should be shown as asterisks")
	}
	
	// Check other fields are shown
	if !strings.Contains(formatted, "john@example.com") {
		t.Error("Email should be visible")
	}
}