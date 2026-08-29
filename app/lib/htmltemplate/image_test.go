package htmltemplate_test

import (
	"testing"

	"github.com/josephspurrier/polarbearblog/app/lib/htmltemplate"
	"github.com/stretchr/testify/assert"
)

func TestFirstImageURL(t *testing.T) {
	tt := []struct {
		name    string
		content string
		output  string
	}{
		{"empty", "", ""},
		{"none", "# Title\n\nJust words.", ""},
		{"markdown", "Words.\n\n![alt](/images/one.png)\n\n![alt](/images/two.png)", "/images/one.png"},
		{"html", `<img src="https://example.com/one.jpg" width="100%">`, "https://example.com/one.jpg"},
		{"html single quotes", `<img alt='a' src='/one.jpg'>`, "/one.jpg"},
		{"html before markdown", "<img src=\"/one.jpg\">\n\n![alt](/two.png)", "/one.jpg"},
		{"markdown before html", "![alt](/one.png)\n\n<img src=\"/two.jpg\">", "/one.png"},
		{"markdown with title", `![alt](/one.png "A title")`, "/one.png"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, htmltemplate.FirstImageURL(tc.content))
		})
	}
}
