package markdown

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// MDPackBasePathKey is the context key for storing the base path
var MDPackBasePathKey = parser.NewContextKey()

// MDPackFile represents an extracted file from mdpack
type MDPackFile struct {
	Name     string
	Content  []byte
	MimeType string
	Encoding string
	Embed    bool
}

// MDPackExtractor interface for file storage
type MDPackExtractor interface {
	AddFile(path string, file *MDPackFile)
	GetFile(path, filename string) (*MDPackFile, bool)
}

// mdpackExtension is a goldmark extension for mdpack support
type mdpackExtension struct {
	extractor MDPackExtractor
	basePath  string
}

// MDPackOption is a functional option for MDPack extension
type MDPackOption func(*mdpackExtension)

// WithMDPackExtractor sets the file extractor
func WithMDPackExtractor(extractor MDPackExtractor) MDPackOption {
	return func(e *mdpackExtension) {
		e.extractor = extractor
	}
}

// WithBasePath sets the base path for file URLs
func WithBasePath(path string) MDPackOption {
	return func(e *mdpackExtension) {
		e.basePath = path
	}
}

// NewMDPackExtension creates a new mdpack extension
func NewMDPackExtension(opts ...MDPackOption) goldmark.Extender {
	e := &mdpackExtension{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *mdpackExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(&mdpackParser{
				extractor: e.extractor,
				basePath:  e.basePath,
			}, 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&mdpackRenderer{
				basePath: e.basePath,
			}, 100),
		),
	)
}

// mdpackParser parses mdpack code blocks
type mdpackParser struct {
	extractor MDPackExtractor
	basePath  string
}

func (p *mdpackParser) Trigger() []byte {
	return []byte{'`'}
}

func (p *mdpackParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if !bytes.HasPrefix(line, []byte("```pack,")) {
		return nil, parser.NoChildren
	}
	
	// Get base path from context
	basePath := p.basePath
	if ctxPath := pc.Get(MDPackBasePathKey); ctxPath != nil {
		if pathStr, ok := ctxPath.(string); ok {
			basePath = pathStr
		}
	}

	// Parse the pack header
	header := string(line[8:])
	parts := strings.SplitN(strings.TrimSpace(header), " ", 2)
	if len(parts) == 0 {
		return nil, parser.NoChildren
	}

	// Parse filename and parameters
	fileSpec := parts[0]
	idx := strings.Index(fileSpec, "?")
	
	var filename string
	params := make(map[string]string)
	
	if idx > 0 {
		filename = fileSpec[:idx]
		paramStr := fileSpec[idx+1:]
		
		// Parse URL parameters
		values, err := url.ParseQuery(paramStr)
		if err == nil {
			for k, v := range values {
				if len(v) > 0 {
					params[k] = v[0]
				}
			}
		}
	} else {
		filename = fileSpec
	}

	// Skip the header line
	reader.Advance(len(line))

	// Read content until closing ```
	var content bytes.Buffer
	for {
		line, _ := reader.PeekLine()
		if bytes.HasPrefix(line, []byte("```")) {
			reader.Advance(len(line))
			break
		}
		if len(line) == 0 {
			return nil, parser.NoChildren
		}
		content.Write(line)
		reader.Advance(len(line))
	}

	// Process the content based on encoding
	fileContent := content.Bytes()
	if encoding := params["encoding"]; encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(fileContent)))
		if err == nil {
			fileContent = decoded
		}
	}

	// Store the file if extractor is available
	if p.extractor != nil {
		file := &MDPackFile{
			Name:     filename,
			Content:  fileContent,
			MimeType: params["mimetype"],
			Encoding: params["encoding"],
			Embed:    params["embed"] == "true",
		}
		p.extractor.AddFile(basePath, file)
		// Debug logging
		// fmt.Printf("MDPack: Stored file %s at path %s (size: %d)\n", filename, basePath, len(fileContent))
	}

	// Create a custom node for rendering
	node := &mdpackNode{
		BaseBlock: ast.BaseBlock{},
		Filename:  filename,
		Params:    params,
		Content:   fileContent,
		BasePath:  basePath,
	}
	return node, parser.NoChildren
}

func (p *mdpackParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

func (p *mdpackParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// Nothing to do
}

func (p *mdpackParser) CanInterruptParagraph() bool {
	return false
}

func (p *mdpackParser) CanAcceptIndentedLine() bool {
	return false
}

// mdpackNode represents an mdpack block in the AST
type mdpackNode struct {
	ast.BaseBlock
	Filename string
	Params   map[string]string
	Content  []byte
	BasePath string
}

// KindMDPack is the NodeKind for MDPack
var KindMDPack = ast.NewNodeKind("MDPack")

func (n *mdpackNode) Kind() ast.NodeKind {
	return KindMDPack
}

func (n *mdpackNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// mdpackRenderer renders mdpack nodes
type mdpackRenderer struct {
	basePath string
}

func (r *mdpackRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMDPack, r.renderMDPack)
}

func (r *mdpackRenderer) renderMDPack(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n, ok := node.(*mdpackNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node type")
	}

	// If embed is true and it's an image, render as img tag
	if n.Params["embed"] == "true" && strings.HasPrefix(n.Params["mimetype"], "image/") {
		encoded := base64.StdEncoding.EncodeToString(n.Content)
		w.WriteString(`<img src="data:`)
		w.WriteString(html.EscapeString(n.Params["mimetype"]))
		w.WriteString(`;base64,`)
		w.WriteString(encoded)
		w.WriteString(`" alt="`)
		_, _ = w.Write([]byte(html.EscapeString(n.Filename)))
		w.WriteString(`" />`)
	} else if n.Params["inline"] == "true" && strings.HasPrefix(n.Params["mimetype"], "image/") {
		// Render inline image with query parameter
		w.WriteString(`<div class="mdpack-inline">`)
		w.WriteString(`<img src="?file=`)
		_, _ = w.Write([]byte(html.EscapeString(n.Filename)))
		w.WriteString(`" alt="`)
		_, _ = w.Write([]byte(html.EscapeString(n.Filename)))
		w.WriteString(`" />`)
		w.WriteString(`</div>`)
	} else if n.Params["iframe"] == "true" && (strings.HasPrefix(n.Params["mimetype"], "text/html") || 
		strings.HasPrefix(n.Params["mimetype"], "text/plain") || 
		strings.HasPrefix(n.Params["mimetype"], "application/pdf")) {
		// Render as iframe with query parameter
		w.WriteString(`<div class="mdpack-iframe">`)
		w.WriteString(`<iframe src="?file=`)
		_, _ = w.Write([]byte(html.EscapeString(n.Filename)))
		w.WriteString(`" title="`)
		_, _ = w.Write([]byte(html.EscapeString(n.Filename)))
		w.WriteString(`"></iframe>`)
		w.WriteString(`</div>`)
	} else {
		// Render as a link with query parameter
		w.WriteString(`<div class="mdpack-file">`)
		w.WriteString(`<a href="?file=`)
		_, _ = w.Write([]byte(html.EscapeString(n.Filename)))
		w.WriteString(`">`)
		_, _ = w.Write([]byte(html.EscapeString(n.Filename)))
		w.WriteString(`</a>`)
		
		// Add file info
		w.WriteString(` <span class="mdpack-info">(`)
		w.WriteString(fmt.Sprintf("%d bytes", len(n.Content)))
		if n.Params["mimetype"] != "" {
			w.WriteString(`, `)
			_, _ = w.Write([]byte(html.EscapeString(n.Params["mimetype"])))
		}
		w.WriteString(`)</span>`)
		w.WriteString(`</div>`)
	}

	return ast.WalkContinue, nil
}