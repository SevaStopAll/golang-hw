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
	done, err := checkRegexp(matched, matched3, matched2)
	if done {
		return "", err
	}
	runes := []rune(str)
	if len(runes) == 0 {
		return "", nil
	}
	for i := 0; i < len(runes); i++ {
		if unicode.IsDigit(runes[i]) && !unicode.IsDigit(runes[i-1]) {
			continue
		}
		if i+1 < len(runes) {
			nextChar := runes[i+1]
			a := string(nextChar)
			b := string(runes[i])
			if strings.EqualFold(b, "\\") {
				continue
			}
			if unicode.IsDigit(nextChar) && !unicode.IsDigit(runes[i]) {
				length, _ := strconv.Atoi(a)
				writeString(length, builder, b)
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

func writeString(length int, builder strings.Builder, b string) {
	for i := 0; i < length; i++ {
		builder.WriteString(b)
	}
}

func checkRegexp(matched bool, matched3 bool, matched2 bool) (bool, error) {
	if (matched && !matched3) || matched2 {
		return true, ErrInvalidString
	}
	return false, nil
}
