package forms

import (
	"strings"
	
	"gno.land/p/demo/ufmt"
	"gno.land/p/moul/md"
)

// State to store form submissions for demonstration
var submissions []Submission

type Submission struct {
	Path   string
	Fields map[string]string
}

func Render(path string) string {
	// Handle different routes
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[0] != "" {
		switch parts[0] {
		case "simple":
			return renderSimpleForm(path)
		case "registration":
			return renderRegistrationForm(path)
		case "survey":
			return renderSurveyForm(path)
		case "contact":
			return renderContactForm(path)
		case "login":
			return renderLoginForm(path)
		case "payment":
			return renderPaymentForm(path)
		case "builder":
			return renderBuilderExample(path)
		case "submissions":
			return renderSubmissions()
		default:
			// Check if it's a form submission result
			if strings.Contains(path, "?") {
				return handleFormSubmission(path)
			}
		}
	}
	
	// Default: show all form examples
	return renderHomepage()
}

func renderHomepage() string {
	return md.H1("Gno Forms Documentation & Examples") +
		md.Paragraph("This realm demonstrates various form patterns and best practices using Gno's form system.") +
		
		md.H2("Available Examples") +
		md.BulletList([]string{
			md.Link("Simple Form", "/r/docs/forms:simple") + " - Basic form with text inputs",
			md.Link("Registration Form", "/r/docs/forms:registration") + " - Complete user registration example",
			md.Link("Survey Form", "/r/docs/forms:survey") + " - Multi-input survey with various field types",
			md.Link("Contact Form", "/r/docs/forms:contact") + " - Standard contact form pattern",
			md.Link("Login Form", "/r/docs/forms:login") + " - Authentication form example",
			md.Link("Payment Form", "/r/docs/forms:payment") + " - Payment details collection",
			md.Link("FormBuilder Example", "/r/docs/forms:builder") + " - Using p/moul/md FormBuilder",
			md.Link("View Submissions", "/r/docs/forms:submissions") + " - See submitted form data (demo only)",
		}) +
		
		md.H2("Quick Start") +
		md.H3("Method 1: Raw HTML Tags") +
		md.CodeBlock(`<gno-form>
  <gno-input name="username" type="text" placeholder="Enter username" />
  <gno-input name="email" type="email" placeholder="Enter email" />
</gno-form>`) +
		
		md.H3("Method 2: Using p/moul/md Helpers") +
		md.LanguageCodeBlock("go", `import "gno.land/p/moul/md"

func Render(path string) string {
    form := md.NewForm("submit").
        AddTextInput("username", "Enter username").
        AddEmailInput("email", "Enter email").
        Build()
    
    return form
}`) +
		
		md.H2("Form Security") +
		md.Paragraph("All forms are processed within the realm where they are defined. This ensures:") +
		md.BulletList([]string{
			"No cross-realm form submissions",
			"Clear visibility of which realm processes the form",
			"Secure handling of user input within the blockchain context",
		}) +
		
		md.H2("Best Practices") +
		md.OrderedList([]string{
			"Always provide clear placeholders for inputs",
			"Use appropriate input types for validation",
			"Handle form submissions gracefully",
			"Provide feedback after form submission",
			"Keep forms simple and focused",
		})
}

func renderSimpleForm(path string) string {
	return md.H1("Simple Form Example") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.Paragraph("This demonstrates a basic form with text inputs.") +
		
		md.H2("Live Example") +
		`<gno-form path="simple/submit">
  <gno-input name="name" type="text" placeholder="Your name" />
  <gno-input name="message" type="text" placeholder="Your message" />
</gno-form>` +
		
		md.H2("Source Code") +
		md.LanguageCodeBlock("markdown", `<gno-form path="simple/submit">
  <gno-input name="name" type="text" placeholder="Your name" />
  <gno-input name="message" type="text" placeholder="Your message" />
</gno-form>`)
}

func renderRegistrationForm(path string) string {
	return md.H1("Registration Form Example") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.Paragraph("A complete user registration form with various input types.") +
		
		md.H2("Live Example") +
		md.NewForm("registration/submit").
			AddTextInput("username", "Choose a username").
			AddEmailInput("email", "Your email address").
			AddPasswordInput("password", "Create a password").
			AddPasswordInput("confirm", "Confirm password").
			AddNumberInput("age", "Your age").
			AddTelInput("phone", "Phone number (optional)").
			Build() +
		
		md.H2("Source Code") +
		md.LanguageCodeBlock("go", `form := md.NewForm("registration/submit").
    AddTextInput("username", "Choose a username").
    AddEmailInput("email", "Your email address").
    AddPasswordInput("password", "Create a password").
    AddPasswordInput("confirm", "Confirm password").
    AddNumberInput("age", "Your age").
    AddTelInput("phone", "Phone number (optional)").
    Build()`)
}

func renderSurveyForm(path string) string {
	return md.H1("Survey Form Example") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.Paragraph("A survey demonstrating different input types and patterns.") +
		
		md.H2("Customer Satisfaction Survey") +
		`<gno-form path="survey/submit">
  <gno-input name="customer_name" type="text" placeholder="Your name" />
  <gno-input name="order_number" type="text" placeholder="Order number" />
  <gno-input name="rating" type="number" placeholder="Rating (1-10)" />
  <gno-input name="email" type="email" placeholder="Email for follow-up" />
  <gno-input name="improvements" type="text" placeholder="Suggestions for improvement" />
</gno-form>`
}

func renderContactForm(path string) string {
	return md.H1("Contact Form Example") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.Paragraph("Standard contact form pattern used across many applications.") +
		
		md.H2("Get in Touch") +
		md.NewForm("contact/submit").
			AddTextInput("name", "Full name").
			AddEmailInput("email", "Email address").
			AddTextInput("subject", "Subject").
			AddTextInput("message", "Your message").
			Build() +
		
		md.H2("Implementation Notes") +
		md.BulletList([]string{
			"Email validation is handled by the browser",
			"Form data is submitted as query parameters",
			"Consider adding rate limiting in production",
			"You might want to integrate with a notification system",
		})
}

func renderLoginForm(path string) string {
	return md.H1("Login Form Example") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.Paragraph("A secure login form example.") +
		
		md.H2("Sign In") +
		`<gno-form path="login/submit">
  <gno-input name="username" type="text" placeholder="Username or email" />
  <gno-input name="password" type="password" placeholder="Password" />
</gno-form>` +
		
		md.H2("Security Considerations") +
		md.BulletList([]string{
			"Passwords are masked using type=\"password\"",
			"Never log or display password values",
			"Consider implementing rate limiting",
			"Add proper authentication logic in your realm",
		})
}

func renderPaymentForm(path string) string {
	return md.H1("Payment Form Example") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.Paragraph("Example of collecting payment details (demonstration only).") +
		
		md.H2("Payment Information") +
		md.NewForm("payment/submit").
			AddTextInput("cardholder", "Cardholder name").
			AddNumberInput("cardnumber", "Card number").
			AddTextInput("expiry", "MM/YY").
			AddNumberInput("cvv", "CVV").
			AddTextInput("billing_address", "Billing address").
			AddNumberInput("amount", "Amount to pay").
			Build() +
		
		md.H2("⚠️ Important Security Note") +
		md.Paragraph("This is for demonstration only. In production:") +
		md.BulletList([]string{
			"Never collect real payment information without proper PCI compliance",
			"Use established payment processors",
			"Implement proper encryption and security measures",
			"Follow all relevant financial regulations",
		})
}

func renderBuilderExample(path string) string {
	return md.H1("FormBuilder Advanced Example") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.Paragraph("This example shows advanced usage of the FormBuilder API.") +
		
		md.H2("Dynamic Form Generation") +
		md.LanguageCodeBlock("go", `// Example: Generate form based on configuration
type FormField struct {
    Name        string
    Type        string
    Placeholder string
}

func generateForm(fields []FormField) string {
    fb := md.NewForm("dynamic/submit")
    
    for _, field := range fields {
        switch field.Type {
        case "email":
            fb.AddEmailInput(field.Name, field.Placeholder)
        case "password":
            fb.AddPasswordInput(field.Name, field.Placeholder)
        case "number":
            fb.AddNumberInput(field.Name, field.Placeholder)
        case "tel":
            fb.AddTelInput(field.Name, field.Placeholder)
        default:
            fb.AddTextInput(field.Name, field.Placeholder)
        }
    }
    
    return fb.Build()
}`) +
		
		md.H2("Live Example - Event Registration") +
		generateEventForm() +
		
		md.H2("Benefits of FormBuilder") +
		md.BulletList([]string{
			"Type-safe form construction",
			"Fluent API for better readability",
			"Easy to generate forms dynamically",
			"Consistent HTML output",
			"Built-in escaping for security",
		})
}

func generateEventForm() string {
	return md.NewForm("event/register").
		AddTextInput("attendee_name", "Full name").
		AddEmailInput("attendee_email", "Email address").
		AddTextInput("organization", "Company/Organization").
		AddTelInput("phone", "Contact number").
		AddNumberInput("guests", "Number of additional guests").
		AddTextInput("dietary", "Dietary requirements").
		AddTextInput("comments", "Additional comments").
		Build()
}

func handleFormSubmission(path string) string {
	// Parse the query parameters
	parts := strings.Split(path, "?")
	if len(parts) < 2 {
		return md.H1("Invalid submission")
	}
	
	// Store submission for demo
	submission := Submission{
		Path:   parts[0],
		Fields: parseQueryParams(parts[1]),
	}
	submissions = append(submissions, submission)
	
	// Generate response
	return md.H1("Form Submitted Successfully! ✅") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.H2("Submission Details") +
		md.Paragraph("Form: " + md.InlineCode(parts[0])) +
		md.H3("Submitted Data:") +
		formatSubmissionData(submission.Fields) +
		md.HorizontalRule() +
		md.Paragraph(md.Link("View all submissions", "/r/docs/forms:submissions")) +
		md.Paragraph(md.Link("Try another form", "/r/docs/forms"))
}

func parseQueryParams(query string) map[string]string {
	params := make(map[string]string)
	pairs := strings.Split(query, "&")
	
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}
	
	return params
}

func formatSubmissionData(fields map[string]string) string {
	var items []string
	for key, value := range fields {
		// Mask password fields
		displayValue := value
		if strings.Contains(strings.ToLower(key), "password") {
			displayValue = "********"
		}
		items = append(items, md.Bold(key) + ": " + md.InlineCode(displayValue))
	}
	return md.BulletList(items)
}

func renderSubmissions() string {
	content := md.H1("Form Submissions (Demo)") +
		md.Link("← Back to Examples", "/r/docs/forms") + "\n\n" +
		md.Paragraph("This page shows all form submissions for demonstration purposes. In a real application, this data would be properly secured.") +
		md.HorizontalRule()
	
	if len(submissions) == 0 {
		content += md.Paragraph("No submissions yet. Try submitting some forms!")
	} else {
		content += md.H2(ufmt.Sprintf("Total Submissions: %d", len(submissions)))
		
		// Show last 10 submissions
		start := 0
		if len(submissions) > 10 {
			start = len(submissions) - 10
		}
		
		for i := len(submissions) - 1; i >= start; i-- {
			sub := submissions[i]
			content += md.H3(ufmt.Sprintf("Submission #%d - Form: %s", i+1, sub.Path)) +
				formatSubmissionData(sub.Fields) +
				md.HorizontalRule()
		}
	}
	
	return content
}