package utils

import (
	"crypto/md5"
	"encoding/hex"
	"html"
	"strings"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

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
