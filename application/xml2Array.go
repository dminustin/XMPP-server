package application

import (
	"regexp"
)

type struc struct {

	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Stream string   `xml:"stream"`
}


func ParseMSGType(s string) (string) {
	var re = regexp.MustCompile("(<\\?xml .*?>)")
	s = re.ReplaceAllString(s, "")
	re = regexp.MustCompile("(<stream)")
	if re.MatchString(s) {
		s = s + "</stream:stream>"
	}
	re = regexp.MustCompile("([a-zA-Z]+)")
	cmd := re.FindString(s)
	return cmd
}