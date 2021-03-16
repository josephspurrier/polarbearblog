package htmltemplate

import (
	"strings"

	"github.com/josephspurrier/polarbearblog/app/model"
	blackfriday "github.com/russross/blackfriday/v2"
	"jaytaylor.com/html2text"
)

// MetaDescription -
func MetaDescription(v model.Post) string {
	unsafeHTML := blackfriday.Run([]byte(v.Content))
	plaintext, err := html2text.FromString(string(unsafeHTML))
	if err != nil {
		plaintext = v.Content
	}
	period := strings.Index(plaintext, ". ")
	if period > 0 {
		plaintext = plaintext[:period+1]
	}

	return plaintext
}
