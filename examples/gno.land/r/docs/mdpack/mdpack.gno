// MDPack demonstrates embedding files in markdown using txtar-like format
package mdpack

func Render(string) string {
	return `# MDPack Example

MDPack allows embedding files directly in markdown content using a txtar-like format.

## Display Modes

### 1. Link Mode (default)
Files are shown as clickable links that open in the browser:

` + "```pack,example.txt --" + `
Hello, World!
This is a text file that will open in your browser when clicked.
` + "```" + `

` + "```pack,data.json?mimetype=application/json --" + `
{
  "name": "MDPack",
  "version": "1.0",
  "features": ["embed", "inline", "base64"]
}
` + "```" + `

### 2. Inline Mode
Images can be displayed inline in the page:

` + "```pack,logo.png?inline=true&encoding=base64&mimetype=image/png --" + `
iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII=
` + "```" + `

### 3. Embed Mode
For data URIs embedded directly in HTML (useful for small images):

` + "```pack,tiny.svg?embed=true&mimetype=image/svg+xml --" + `
<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20">
  <circle cx="10" cy="10" r="8" fill="#0969da"/>
</svg>
` + "```" + `

## Parameters

- **encoding**: Specify encoding (e.g., base64)
- **mimetype**: Set MIME type for proper display
- **embed**: Embed as data URI (true/false)
- **inline**: Display images inline (true/false)

## Use Cases

1. **Documentation**: Include example files, configs, and images
2. **Tutorials**: Embed code samples with proper syntax highlighting
3. **Demos**: Show images and data files without external hosting
4. **Self-contained content**: Everything in one markdown file

## Magic URLs 🪄

The best part? You can use standard markdown image syntax!

` + "```pack,magic.png?encoding=base64&mimetype=image/png --" + `
iVBORw0KGgoAAAANSUhEUgAAAAgAAAAIAQMAAAD+wSzIAAAABlBMVEX///+/v7+jQ3Y5AAAADklEQVQI12P4AIX8EAgALgAD/aNpbtEAAAAASUVORK5CYII=
` + "```" + `

Now you can reference it with regular markdown: ![Magic](?file=magic.png)

The file is extracted during rendering and served as if it were a real file!

## Caching

Files are cached for only 30 seconds by default, ensuring fresh content while still
providing the "magic" experience. This short cache time means:

- Files are available immediately after rendering
- Memory usage is minimal
- Content updates are reflected quickly
- The feature can be easily disabled if needed`
}
