package md_test

import (
	"strings"
	"testing"

	"gno.land/p/moul/md"
)

func TestForm(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		content  string
		expected string
	}{
		{
			name:    "form with path",
			path:    "user",
			content: "test content",
			expected: "<gno-form path=\"user\">\ntest content\n</gno-form>\n",
		},
		{
			name:    "form without path",
			path:    "",
			content: "test content",
			expected: "<gno-form>\ntest content\n</gno-form>\n",
		},
		{
			name:    "form with content ending in newline",
			path:    "submit",
			content: "test content\n",
			expected: "<gno-form path=\"submit\">\ntest content\n</gno-form>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := md.Form(tt.path, tt.content)
			if result != tt.expected {
				t.Errorf("Form() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormInput(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		inputType   string
		placeholder string
		expected    string
		shouldPanic bool
	}{
		{
			name:        "basic text input",
			inputName:   "username",
			inputType:   "text",
			placeholder: "Enter username",
			expected:    "  <gno-input name=\"username\" type=\"text\" placeholder=\"Enter username\" />\n",
		},
		{
			name:        "number input",
			inputName:   "age",
			inputType:   "number",
			placeholder: "Your age",
			expected:    "  <gno-input name=\"age\" type=\"number\" placeholder=\"Your age\" />\n",
		},
		{
			name:        "default type and placeholder",
			inputName:   "field",
			inputType:   "",
			placeholder: "",
			expected:    "  <gno-input name=\"field\" type=\"text\" placeholder=\"Enter value\" />\n",
		},
		{
			name:        "empty name should panic",
			inputName:   "",
			inputType:   "text",
			placeholder: "test",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("FormInput() did not panic with empty name")
					}
				}()
			}
			
			result := md.FormInput(tt.inputName, tt.inputType, tt.placeholder)
			if !tt.shouldPanic && result != tt.expected {
				t.Errorf("FormInput() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormBuilder(t *testing.T) {
	t.Run("complete form with multiple inputs", func(t *testing.T) {
		form := md.NewForm("register").
			AddTextInput("username", "Choose a username").
			AddEmailInput("email", "Your email address").
			AddPasswordInput("password", "Create a password").
			AddNumberInput("age", "Your age").
			AddTelInput("phone", "Phone number").
			Build()

		// Check that the form starts and ends correctly
		if !strings.HasPrefix(form, "<gno-form path=\"register\">\n") {
			t.Errorf("Form should start with <gno-form path=\"register\">")
		}
		if !strings.HasSuffix(form, "</gno-form>\n") {
			t.Errorf("Form should end with </gno-form>")
		}

		// Check for each input
		expectedInputs := []string{
			`<gno-input name="username" type="text" placeholder="Choose a username" />`,
			`<gno-input name="email" type="email" placeholder="Your email address" />`,
			`<gno-input name="password" type="password" placeholder="Create a password" />`,
			`<gno-input name="age" type="number" placeholder="Your age" />`,
			`<gno-input name="phone" type="tel" placeholder="Phone number" />`,
		}

		for _, expected := range expectedInputs {
			if !strings.Contains(form, expected) {
				t.Errorf("Form should contain %q", expected)
			}
		}
	})

	t.Run("form without path", func(t *testing.T) {
		form := md.NewForm("").
			AddTextInput("search", "Search...").
			Build()

		if !strings.HasPrefix(form, "<gno-form>\n") {
			t.Errorf("Form without path should start with <gno-form>")
		}
	})
}