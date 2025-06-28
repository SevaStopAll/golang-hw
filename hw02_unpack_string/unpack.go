package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	builder := strings.Builder{}
	matched3, _ := regexp.MatchString(`\\\d{2}`, str)
	matched, _ := regexp.MatchString("\\d{2}", str)
	matched2, _ := regexp.MatchString("\\A\\d", str)
	if (matched && !matched3) || matched2 {
		return "", ErrInvalidString
	}
	runes := []rune(str)
	if len(runes) == 0 {
		return "", nil
	}
	for i := 0; i < len(runes); i++ {
		if unicode.IsDigit(runes[i]) && !unicode.IsDigit(runes[i-1]) {
			continue
		} else if strings.EqualFold(string(runes[i]), "\\") {
			continue
		}
		if i+1 < len(runes) {
			nextChar := runes[i+1]
			if unicode.IsDigit(nextChar) && !unicode.IsDigit(runes[i]) {
				a := string(nextChar)
				b := string(runes[i])
				length, _ := strconv.Atoi(a)
				for i := 0; i < length; i++ {
					builder.WriteString(b)
				}
			} else {
				a := string(runes[i])
				builder.WriteString(a)
			}
		} else if !unicode.IsDigit(runes[i]) {
			a := string(runes[i])
			builder.WriteString(a)
		}
	}
	return builder.String(), nil
}
