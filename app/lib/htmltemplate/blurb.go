package htmltemplate

import (
	"strings"

	blackfriday "github.com/russross/blackfriday/v2"
	"jaytaylor.com/html2text"
)

// PlaintextBlurb returns a plaintext blurb from markdown content.
func PlaintextBlurb(s string) string {
	unsafeHTML := blackfriday.Run([]byte(s))
	plaintext, err := html2text.FromString(string(unsafeHTML))
	if err != nil {
		plaintext = s
	}
	period := strings.Index(plaintext, ". ")
	if period > 0 {
		plaintext = plaintext[:period+1]
	}

	return plaintext
}
