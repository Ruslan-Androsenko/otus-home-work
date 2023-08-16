package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const backSlash rune = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	var hasPrevRuneShielded, hasPrevBackSlashWritten bool
	inputRunes := []rune(inputString)
	builder := strings.Builder{}

	if len(inputRunes) > 0 && unicode.IsDigit(inputRunes[0]) {
		return "", ErrInvalidString
	}

	for index, item := range inputRunes {
		var prevRune, nextRune rune

		if index > 0 {
			prevRune = inputRunes[index-1]
		}

		if index+1 < len(inputRunes) {
			nextRune = inputRunes[index+1]
		}

		if index > 1 && unicode.IsDigit(item) && unicode.IsDigit(prevRune) && !hasPrevRuneShielded {
			return "", ErrInvalidString
		}

		if hasEscapedNotNumber(prevRune, item) || hasIncorrectEscaping(prevRune, item, nextRune) {
			return "", ErrInvalidString
		}

		hasWriteBackSlashCharacter := hasWriteBackSlash(prevRune, item, nextRune)

		switch {
		// Необходимо записать обратный слеш и установить флаг записи.
		case hasWriteBackSlashCharacter && !hasPrevBackSlashWritten:
			builder.WriteString(string(item))
			hasPrevBackSlashWritten = true

		// На предыдущей итерации обратный слеш был записан, необходимо снять этот флаг записи.
		case hasWriteBackSlashCharacter && hasPrevBackSlashWritten:
			hasPrevBackSlashWritten = false

		case hasSetShieldedFlag(prevRune, item) && !hasPrevBackSlashWritten:
			hasPrevRuneShielded = true

		case hasWriteShieldedNumber(prevRune, item) && !hasPrevRuneShielded:
			if !unicode.IsDigit(nextRune) {
				builder.WriteString(string(item))
			} else {
				hasPrevRuneShielded = true
			}

		// Необходимо ли записать символ с указанным повторением.
		case unicode.IsDigit(item) || hasPrevRuneShielded:
			countSymbols, _ := strconv.Atoi(string(item))
			builder.WriteString(strings.Repeat(string(prevRune), countSymbols))
			hasPrevRuneShielded = false

		case hasRegularCharacter(nextRune, item):
			builder.WriteString(string(item))
		}
	}

	return builder.String(), nil
}

// Экранировано не число.
func hasEscapedNotNumber(prevRune, item rune) bool {
	return prevRune == backSlash && item != backSlash && !unicode.IsDigit(item)
}

// Неверное экранирование, т.к. обратный слэш находится в конце строки.
func hasIncorrectEscaping(prevRune, item, nextRune rune) bool {
	return prevRune != backSlash && item == backSlash && nextRune == 0
}

// Необходимо ли установит флаг экраннированного символа.
func hasSetShieldedFlag(prevRune, item rune) bool {
	return prevRune == backSlash && item == backSlash
}

// Необходимо ли записать обратный слэш.
func hasWriteBackSlash(prevRune, item, nextRune rune) bool {
	return prevRune == backSlash && item == backSlash && (nextRune == backSlash || nextRune == 0)
}

// Необходимо ли записать экранированное число.
func hasWriteShieldedNumber(prevRune, item rune) bool {
	return prevRune == backSlash && unicode.IsDigit(item)
}

// Необходимо ли записать обычный символ.
func hasRegularCharacter(nextRune, item rune) bool {
	return nextRune == 0 || !unicode.IsDigit(nextRune) && item != backSlash
}
