package htmltemplate

import (
	"fmt"
	"html/template"
	"strings"

	blackfriday "github.com/russross/blackfriday/v2"
)

func sanitizedContent(t *template.Template, content string) (*template.Template, error) {
	b := []byte(content)
	// Ensure unit line endings are used when pulling out of JSON.
	markdownWithUnixLineEndings := strings.Replace(string(b), "\r\n", "\n", -1)
	unsafeHTML := blackfriday.Run([]byte(markdownWithUnixLineEndings))
	//safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)
	safeContent := fmt.Sprintf(`[{[{define "content"}]}]%s[{[{end}]}]`, string(unsafeHTML))
	t = t.Delims("[{[{", "}]}]") // Change delimiters temporarily so code samples can use Go blocks.
	var err error
	t, err = t.Parse(safeContent)
	if err != nil {
		return nil, err
	}
	t = t.Delims("{{", "}}") // Reset delimiters
	return t, nil
}
