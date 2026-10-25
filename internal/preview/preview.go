package preview

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
)

const defaultMaxRunes = 100

func ContentPreview(raw string, maxRunes int) string {
	if maxRunes <= 0 {
		maxRunes = defaultMaxRunes
	}

	plain := extractText(raw)
	plain = strings.Join(strings.Fields(plain), " ")

	if utf8.RuneCountInString(plain) <= maxRunes {
		return plain
	}

	return string([]rune(plain)[:maxRunes])
}

func extractText(raw string) string {
	var b strings.Builder
	z := html.NewTokenizer(strings.NewReader(raw))
	skipDepth := 0

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return b.String()
		case html.StartTagToken:
			if isSkippedTag(z) {
				skipDepth++
			}
		case html.EndTagToken:
			if isSkippedTag(z) && skipDepth > 0 {
				skipDepth--
			}
		case html.TextToken:
			if skipDepth == 0 {
				b.Write(z.Text())
			}
		}
	}
}

func isSkippedTag(z *html.Tokenizer) bool {
	name, _ := z.TagName()
	return bytes.Equal(name, []byte("script")) ||
		bytes.Equal(name, []byte("style")) ||
		bytes.Equal(name, []byte("noscript"))
}
