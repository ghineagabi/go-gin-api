package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func OnlyUnicode(s string) bool {
	if s = strings.TrimSpace(s); s == "" {
		return false
	}
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func ValidatePassword(s string) bool {
	letters := 0
	var number, upper, sevenOrMore bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsLetter(c) || c == ' ':
			letters++
		}
	}
	sevenOrMore = letters >= 7 && letters < 49
	if sevenOrMore && number && upper {
		return true
	}
	return false
}

func ValidPhoneNumber(s string) bool {
	var validID = regexp.MustCompile(`^(\+4|)?(07[0-8]{1}[0-9]{1}|02[0-9]{2}|03[0-9]{2}){1}?(\s|\.|\-)?([0-9]{3}(\s|\.|\-|)){2}$`)

	return validID.MatchString(s)
}
