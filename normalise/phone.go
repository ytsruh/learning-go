package normalise

import (
	"bytes"
	"regexp"
)

func PhoneNumber(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func PhoneRegex(phone string) string {
	re := regexp.MustCompile("[^0-9]")
	//re := regexp.MustCompile(`\\D`)
	return re.ReplaceAllString(phone, "")
}
