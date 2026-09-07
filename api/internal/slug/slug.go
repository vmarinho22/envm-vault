package slug

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// FromName turns a display name into a URL-safe slug: lowercase ASCII,
// accents stripped, non-alphanumeric runs collapsed to a single hyphen.
func FromName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	deaccented, _, err := transform.String(t, trimmed)
	if err != nil {
		deaccented = trimmed
	}

	return strings.Trim(nonAlnum.ReplaceAllString(strings.ToLower(deaccented), "-"), "-")
}
