package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var b strings.Builder
	var br rune = 0
	var bbr rune = 0
	for _, r := range str {

		if bbr != '\\' && br == '\\' && (r == '\\' || unicode.IsDigit(r)) {

		} else if bbr == '\\' && br == '\\' && r == '\\' {
			b.WriteRune(br)
			br = 0
		} else if unicode.IsDigit(r) {
			if br == 0 {
				return "", ErrInvalidString
			}
			for i := 0; i < int(r)-'0'; i++ {
				b.WriteRune(br)
			}
			r = 0
		} else if br != 0 {
			b.WriteRune(br)
		}
		bbr = br
		br = r
	}
	if br != 0 {
		b.WriteRune(br)
	}
	return b.String(), nil
}
