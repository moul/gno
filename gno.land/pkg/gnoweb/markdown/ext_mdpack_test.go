package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// mockExtractor implements MDPackExtractor for testing
type mockExtractor struct {
	files map[string]*MDPackFile
}

func newMockExtractor() *mockExtractor {
	return &mockExtractor{
		files: make(map[string]*MDPackFile),
	}
}

func (m *mockExtractor) AddFile(path string, file *MDPackFile) {
	key := path + "/" + file.Name
	m.files[key] = file
}

func (m *mockExtractor) GetFile(path, filename string) (*MDPackFile, bool) {
	key := path + "/" + filename
	file, ok := m.files[key]
	return file, ok
}

func TestMDPackExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantHTML string
		wantFile bool
		fileName string
		fileSize int
	}{
		{
			name: "simple text file",
			input: `# Test Document

Here's a file:

` + "```pack,test.txt" + ` --
Hello, World!
This is a test file.
` + "```" + `
`,
			wantHTML: `<h1>Test Document</h1>
<p>Here's a file:</p>
<div class="mdpack-file"><a href="?file=test.txt">test.txt</a> <span class="mdpack-info">(35 bytes)</span></div>`,
			wantFile: true,
			fileName: "test.txt",
			fileSize: 35,
		},
		{
			name: "base64 encoded image",
			input: `# Image Test

` + "```pack,logo.png?encoding=base64&mimetype=image/png" + ` --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
` + "```" + `
`,
			wantHTML: `<h1>Image Test</h1>
<div class="mdpack-file"><a href="?file=logo.png">logo.png</a> <span class="mdpack-info">(70 bytes, image/png)</span></div>`,
			wantFile: true,
			fileName: "logo.png",
			fileSize: 70,
		},
		{
			name: "embedded image",
			input: `# Embedded Image

` + "```pack,inline.svg?embed=true&mimetype=image/svg%2Bxml" + ` --
<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
  <circle cx="50" cy="50" r="40" fill="red" />
</svg>
` + "```" + `
`,
			wantHTML: `<h1>Embedded Image</h1>
<img src="data:image/svg+xml;base64,` + "PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxMDAiIGhlaWdodD0iMTAwIj4KICA8Y2lyY2xlIGN4PSI1MCIgY3k9IjUwIiByPSI0MCIgZmlsbD0icmVkIiAvPgo8L3N2Zz4K" + `" alt="inline.svg" />`,
			wantFile: true,
			fileName: "inline.svg",
			fileSize: 120,
		},
		{
			name: "inline image",
			input: `# Inline Image

` + "```pack,inline.png?inline=true&encoding=base64&mimetype=image/png" + ` --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
` + "```" + `
`,
			wantHTML: `<h1>Inline Image</h1>
<div class="mdpack-inline"><img src="?file=inline.png" alt="inline.png" /></div>`,
			wantFile: true,
			fileName: "inline.png",
			fileSize: 70,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock extractor
			extractor := newMockExtractor()

			// Create goldmark with MDPack extension
			md := goldmark.New(
				goldmark.WithExtensions(
					NewMDPackExtension(
						WithMDPackExtractor(extractor),
						WithBasePath("/test"),
					),
				),
			)

			// Parse and render
			ctx := parser.NewContext()
			ctx.Set(MDPackBasePathKey, "/test")
			
			var buf bytes.Buffer
			doc := md.Parser().Parse(text.NewReader([]byte(tt.input)), parser.WithContext(ctx))
			if err := md.Renderer().Render(&buf, []byte(tt.input), doc); err != nil {
				t.Fatalf("failed to render: %v", err)
			}

			// Check HTML output
			got := buf.String()
			if got != tt.wantHTML {
				t.Errorf("HTML mismatch\ngot:\n%s\nwant:\n%s", got, tt.wantHTML)
			}

			// Check file extraction
			if tt.wantFile {
				file, ok := extractor.GetFile("/test", tt.fileName)
				if !ok {
					t.Errorf("expected file %s to be extracted", tt.fileName)
				} else {
					if len(file.Content) != tt.fileSize {
						t.Errorf("file size mismatch: got %d, want %d", len(file.Content), tt.fileSize)
					}
				}
			}
		})
	}
}

func TestMDPackParser(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantFilename  string
		wantMimeType  string
		wantEncoding  string
		wantEmbed     bool
		wantContent   string
	}{
		{
			name:         "simple file",
			input:        "```pack,test.txt --\nHello\n```",
			wantFilename: "test.txt",
			wantContent:  "Hello\n",
		},
		{
			name:         "file with parameters",
			input:        "```pack,data.json?mimetype=application/json --\n{\"test\": true}\n```",
			wantFilename: "data.json",
			wantMimeType: "application/json",
			wantContent:  "{\"test\": true}\n",
		},
		{
			name:         "base64 encoded",
			input:        "```pack,encoded.bin?encoding=base64 --\nSGVsbG8gV29ybGQ=\n```",
			wantFilename: "encoded.bin",
			wantEncoding: "base64",
			wantContent:  "Hello World",
		},
		{
			name:         "embedded file",
			input:        "```pack,embed.txt?embed=true --\nEmbedded content\n```",
			wantFilename: "embed.txt",
			wantEmbed:    true,
			wantContent:  "Embedded content\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := newMockExtractor()
			
			// Create goldmark with MDPack extension to test the parser
			md := goldmark.New(
				goldmark.WithExtensions(
					NewMDPackExtension(
						WithMDPackExtractor(extractor),
						WithBasePath("/test"),
					),
				),
			)
			
			// Parse and render the markdown
			var buf bytes.Buffer
			ctx := parser.NewContext()
			ctx.Set(MDPackBasePathKey, "/test")
			doc := md.Parser().Parse(text.NewReader([]byte(tt.input)), parser.WithContext(ctx))
			if err := md.Renderer().Render(&buf, []byte(tt.input), doc); err != nil {
				t.Fatalf("failed to render: %v", err)
			}
			
			// Check that the file was extracted with correct properties
			file, ok := extractor.GetFile("/test", tt.wantFilename)
			if !ok {
				t.Fatalf("expected file %q to be extracted", tt.wantFilename)
			}

			if string(file.Content) != tt.wantContent {
				t.Errorf("content mismatch: got %q, want %q", string(file.Content), tt.wantContent)
			}

			if file.MimeType != tt.wantMimeType {
				t.Errorf("mimetype mismatch: got %q, want %q", file.MimeType, tt.wantMimeType)
			}

			if file.Encoding != tt.wantEncoding {
				t.Errorf("encoding mismatch: got %q, want %q", file.Encoding, tt.wantEncoding)
			}

			if file.Embed != tt.wantEmbed {
				t.Errorf("embed mismatch: got %v, want %v", file.Embed, tt.wantEmbed)
			}
		})
	}
}

func TestMDPackStorage(t *testing.T) {
	// Test that non-pack code blocks are not affected
	input := `# Regular Code

` + "```go" + `
package main

func main() {
    println("Hello")
}
` + "```" + `

` + "```pack,test.txt" + ` --
This is mdpack
` + "```" + `
`

	extractor := newMockExtractor()
	md := goldmark.New(
		goldmark.WithExtensions(
			NewMDPackExtension(WithMDPackExtractor(extractor)),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("failed to convert: %v", err)
	}

	// Check that Go code block is preserved
	output := buf.String()
	if !strings.Contains(output, "<pre><code class=\"language-go\">") {
		t.Error("expected Go code block to be preserved")
	}

	// Check that mdpack was processed
	if !strings.Contains(output, "mdpack-file") {
		t.Error("expected mdpack file to be processed")
	}

	// Check file extraction
	if _, ok := extractor.GetFile("", "test.txt"); !ok {
		t.Error("expected test.txt to be extracted")
	}
}