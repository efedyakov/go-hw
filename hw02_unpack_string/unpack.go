package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var b strings.Builder
	var br rune
	var bbr rune
	for _, r := range str {
		switch {
		case bbr != '\\' && unicode.IsDigit(br) && unicode.IsDigit(r):
			{
				return "", ErrInvalidString
			}
		case bbr != '\\' && br == '\\' && (r == '\\' || unicode.IsDigit(r)):
			{
			}
		case bbr == '\\' && br == '\\' && r == '\\':
			{
				b.WriteRune(br)
				br = 0
			}
		case unicode.IsDigit(r):
			{
				if br == 0 {
					return "", ErrInvalidString
				}
				for i := 0; i < int(r)-'0'; i++ {
					//if bbr == '\\' {
					//	b.WriteRune(bbr)
					//}
					b.WriteRune(br)
				}
				r = 0
			}
		case br != 0:
			{
				b.WriteRune(br)
			}
		}
		bbr = br
		br = r
	}
	if br != 0 {
		b.WriteRune(br)
	}
	if bbr != '\\' && br == '\\' {
		return "", ErrInvalidString
	}
	return b.String(), nil
}
