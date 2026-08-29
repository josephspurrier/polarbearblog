package htmltemplate

import "regexp"

var (
	htmlImage     = regexp.MustCompile(`(?i)<img[^>]+src=["']([^"']+)["']`)
	markdownImage = regexp.MustCompile(`!\[[^\]]*\]\(\s*<?([^\s)>]+)`)
)

// FirstImageURL returns the URL of the first image in the content, whether it
// is written as an HTML tag or as markdown. It returns an empty string when
// the content contains no images.
func FirstImageURL(s string) string {
	url := ""
	first := -1

	for _, re := range []*regexp.Regexp{htmlImage, markdownImage} {
		m := re.FindStringSubmatchIndex(s)
		if m == nil {
			continue
		}

		if first < 0 || m[0] < first {
			first = m[0]
			url = s[m[2]:m[3]]
		}
	}

	return url
}
