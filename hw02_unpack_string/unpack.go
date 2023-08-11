package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const backSlash rune = '\\'

var ErrInvalidString = errors.New("invalid string")

var hasPrevRuneShielded, hasPrevBackSlashWritten bool

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

		if index > 1 && unicode.IsDigit(item) && unicode.IsDigit(prevRune) && !hasPrevRuneShielded {
			return "", ErrInvalidString
		}

		if prevRune == backSlash && item != backSlash && !unicode.IsDigit(item) {
			return "", ErrInvalidString
		}

		switch {
		case hasWriteBackSlash(prevRune, item, nextRune):
			builder.WriteString(string(item))
			hasPrevBackSlashWritten = true

		case hasSetShieldedFlag(prevRune, item):
			hasPrevRuneShielded = true

		case hasWriteShieldedNumber(prevRune, item):
			if !unicode.IsDigit(nextRune) {
				builder.WriteString(string(item))
			} else {
				hasPrevRuneShielded = true
			}

		case hasWriteRepeatedCharacter(item):
			countSymbols, _ := strconv.Atoi(string(item))
			builder.WriteString(strings.Repeat(string(prevRune), countSymbols))
			hasPrevRuneShielded = false

		case hasRegularCharacter(nextRune, item):
			builder.WriteString(string(item))
		}
	}

	return builder.String(), nil
}

// Необходимо ли установит флаг экраннированного символа.
func hasSetShieldedFlag(prevRune, item rune) bool {
	return prevRune == backSlash && item == backSlash && !hasPrevBackSlashWritten
}

// Необходимо ли записать обратный слэш.
func hasWriteBackSlash(prevRune, item, nextRune rune) bool {
	return prevRune == backSlash && item == backSlash && nextRune == backSlash
}

// Необходимо ли записать экранированное число.
func hasWriteShieldedNumber(prevRune, item rune) bool {
	return prevRune == backSlash && unicode.IsDigit(item) && !hasPrevRuneShielded
}

// Необходимо ли записать символ с указанным повторением.
func hasWriteRepeatedCharacter(item rune) bool {
	return unicode.IsDigit(item) || hasPrevRuneShielded
}

// Необходимо ли записать обычный символ.
func hasRegularCharacter(nextRune, item rune) bool {
	return nextRune == 0 || !unicode.IsDigit(nextRune) && item != backSlash
}
