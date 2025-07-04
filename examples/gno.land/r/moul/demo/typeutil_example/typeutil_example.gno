// Package typeutil_example demonstrates the use of p/moul/typeutil utilities
package typeutil_example

import (
	"strings"

	"gno.land/p/demo/ufmt"
	"gno.land/p/moul/md"
	"gno.land/p/moul/mdtable"
	"gno.land/p/moul/typeutil"
)

// Storage for configuration-like data
var configData = map[string]string{
	"max_users":       "100",
	"enable_feature":  "true",
	"timeout_seconds": "30",
	"ratio":           "0.85",
	"empty_value":     "",
	"invalid_number":  "not-a-number",
	"boolean_yes":     "yes",
	"boolean_1":       "1",
	"boolean_false":   "false",
}

func Render(path string) string {
	parts := strings.Split(path, "/")
	
	if len(parts) == 0 || parts[0] == "" {
		return renderHome()
	}
	
	switch parts[0] {
	case "conversions":
		return renderConversions()
	case "validation":
		return renderValidation()
	case "safe-parsing":
		return renderSafeParsing()
	default:
		return "Unknown example"
	}
}

func renderHome() string {
	out := md.H1("Type Utilities Example with p/moul/typeutil") + "\n\n"
	out += md.P("This realm demonstrates type conversion and validation utilities.") + "\n\n"
	
	out += md.H2("Available Examples") + "\n"
	examples := []string{
		md.BulletItem(md.Link("Type Conversions - Convert between types", "/r/demo/typeutil_example:conversions")),
		md.BulletItem(md.Link("Validation - Check if values are zero/empty", "/r/demo/typeutil_example:validation")),
		md.BulletItem(md.Link("Safe Parsing - Parse config values safely", "/r/demo/typeutil_example:safe-parsing")),
	}
	out += strings.Join(examples, "\n") + "\n\n"
	
	out += md.H2("Sample Configuration Data") + "\n"
	table := mdtable.New([]string{"Key", "Value"})
	for k, v := range configData {
		displayValue := v
		if v == "" {
			displayValue = md.Italic("(empty)")
		}
		table.AddRow([]string{md.Code(k), displayValue})
	}
	out += table.String()
	
	return out
}

func renderConversions() string {
	out := md.H2("Type Conversions Example") + "\n\n"
	out += md.P("Demonstrating various type conversion functions:") + "\n\n"
	
	// Test values
	testValues := []interface{}{
		42,
		"123",
		true,
		false,
		3.14,
		"hello",
		"",
		nil,
	}
	
	table := mdtable.New([]string{"Original Value", "Type", "ToString", "ToInt", "ToBool"})
	
	for _, val := range testValues {
		typeStr := ufmt.Sprintf("%T", val)
		strVal := typeutil.ToString(val)
		intVal := typeutil.ToInt(val, -1) // -1 as default
		boolVal := typeutil.ToBool(val, false)
		
		// Format display values
		displayVal := ufmt.Sprintf("%v", val)
		if val == nil {
			displayVal = "nil"
		} else if val == "" {
			displayVal = `""`
		}
		
		table.AddRow([]string{
			displayVal,
			typeStr,
			ufmt.Sprintf(`"%s"`, strVal),
			ufmt.Sprintf("%d", intVal),
			ufmt.Sprintf("%t", boolVal),
		})
	}
	
	out += table.String() + "\n\n"
	out += md.P(md.Italic("Note: Default values are used when conversion fails"))
	
	return out
}

func renderValidation() string {
	out := md.H2("Validation Example - Zero Value Checking") + "\n\n"
	out += md.P("Using typeutil.IsZero to check for zero/empty values:") + "\n\n"
	
	testCases := []struct {
		name  string
		value interface{}
	}{
		{"Empty string", ""},
		{"Non-empty string", "hello"},
		{"Zero int", 0},
		{"Non-zero int", 42},
		{"False bool", false},
		{"True bool", true},
		{"Nil value", nil},
		{"Empty slice", []string{}},
		{"Non-empty slice", []string{"a", "b"}},
	}
	
	table := mdtable.New([]string{"Test Case", "Value", "IsZero Result"})
	
	for _, tc := range testCases {
		isZero := typeutil.IsZero(tc.value)
		
		// Format value for display
		displayVal := ufmt.Sprintf("%v", tc.value)
		if tc.value == nil {
			displayVal = "nil"
		} else if tc.value == "" {
			displayVal = `""`
		}
		
		resultStr := "No"
		if isZero {
			resultStr = "Yes"
		}
		
		table.AddRow([]string{tc.name, displayVal, resultStr})
	}
	
	out += table.String()
	
	return out
}

func renderSafeParsing() string {
	out := md.H2("Safe Parsing Example - Configuration Values") + "\n\n"
	out += md.P("Safely parsing configuration values with defaults:") + "\n\n"
	
	// Parse configuration with safe defaults
	results := []struct {
		key         string
		description string
		parsed      string
	}{
		{
			key:         "max_users",
			description: "Parse as integer with default 50",
			parsed:      ufmt.Sprintf("%d", typeutil.ToInt(configData["max_users"], 50)),
		},
		{
			key:         "enable_feature",
			description: "Parse as boolean with default false",
			parsed:      ufmt.Sprintf("%t", typeutil.ToBool(configData["enable_feature"], false)),
		},
		{
			key:         "timeout_seconds",
			description: "Parse as integer with default 60",
			parsed:      ufmt.Sprintf("%d", typeutil.ToInt(configData["timeout_seconds"], 60)),
		},
		{
			key:         "invalid_number",
			description: "Parse invalid as integer with default 0",
			parsed:      ufmt.Sprintf("%d", typeutil.ToInt(configData["invalid_number"], 0)),
		},
		{
			key:         "empty_value",
			description: "Parse empty as integer with default 100",
			parsed:      ufmt.Sprintf("%d", typeutil.ToInt(configData["empty_value"], 100)),
		},
		{
			key:         "boolean_yes",
			description: "Parse 'yes' as boolean",
			parsed:      ufmt.Sprintf("%t", typeutil.ToBool(configData["boolean_yes"], false)),
		},
		{
			key:         "boolean_1",
			description: "Parse '1' as boolean",
			parsed:      ufmt.Sprintf("%t", typeutil.ToBool(configData["boolean_1"], false)),
		},
	}
	
	table := mdtable.New([]string{"Config Key", "Raw Value", "Description", "Parsed Result"})
	
	for _, r := range results {
		rawVal := configData[r.key]
		if rawVal == "" {
			rawVal = md.Italic("(empty)")
		}
		
		table.AddRow([]string{
			md.Code(r.key),
			rawVal,
			r.description,
			md.Bold(r.parsed),
		})
	}
	
	out += table.String() + "\n\n"
	
	out += md.H3("Benefits of Safe Parsing") + "\n"
	out += md.BulletItem("Never panics on invalid input") + "\n"
	out += md.BulletItem("Always returns a usable value") + "\n"
	out += md.BulletItem("Explicit defaults make code self-documenting") + "\n"
	out += md.BulletItem("Perfect for parsing user input or configuration") + "\n"
	
	return out
}