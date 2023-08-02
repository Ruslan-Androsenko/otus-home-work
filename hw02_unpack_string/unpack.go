package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	inputRunes := []rune(inputString)
	builder := strings.Builder{}

	for index, item := range inputRunes {
		var prevRune, nextRune rune

		if index == 0 && unicode.IsDigit(item) {
			return "", ErrInvalidString
		}

		if index > 0 {
			prevRune = inputRunes[index-1]
		}

		if index+1 < len(inputRunes) {
			nextRune = inputRunes[index+1]
		}

		if index > 1 && unicode.IsDigit(item) && unicode.IsDigit(prevRune) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(item) {
			countSymbols, _ := strconv.Atoi(string(item))
			builder.WriteString(strings.Repeat(string(prevRune), countSymbols))
		} else if nextRune == 0 || !unicode.IsDigit(nextRune) {
			builder.WriteString(string(item))
		}
	}

	return builder.String(), nil
}
