package utils

import (
	"html"
	"strings"
)

func ReplaceTemplate(s string, template map[string]string) string {

	for k, v := range template {
		s = strings.Replace(s, k, v, 99999)
	}
	return s
}

func QuoteText(s string) string {
	s = html.EscapeString(strings.Trim(s, " "))
	return s
}

func UnQuoteText(s string) string {
	s = html.UnescapeString(strings.Trim(s, " "))
	return s
}
