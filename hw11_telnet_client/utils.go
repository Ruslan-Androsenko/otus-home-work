package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Проверяем корректность входных параметров.
func hasCorrectArguments(args []string) bool {
	countArgs := len(args)

	if n, ok := hasExistsFlags(args); ok {
		countArgs -= n
	}

	return countArgs >= 3
}

// Проверяем наличие переданного флага.
func hasExistsFlags(args []string) (n int, hasFound bool) {
	for _, argItem := range args {
		if strings.HasPrefix(argItem, "-") {
			n++
			hasFound = true
		}
	}

	return
}

// Проверяем что это пришла ошибка а не конец файла.
func hasBeenError(err error) bool {
	return err != nil && !errors.Is(err, io.EOF)
}

// Остановить горутины если ошибка является концом файла.
func stopGoroutinesIfEndOfFile(err error, stop context.CancelFunc) {
	if errors.Is(err, io.EOF) {
		stop()
	}
}

// Прочитать сообщение из входного источника, и записать его в выходной.
func readAndWrite(in io.ReadCloser, out io.Writer) error {
	buffer := make([]byte, 1024)
	n, errRead := in.Read(buffer)
	if hasBeenError(errRead) {
		return fmt.Errorf("cannot read from input: %w", errRead)
	}

	_, errWrite := out.Write(buffer[:n])
	if errWrite != nil {
		return fmt.Errorf("cannot write to output: %w", errWrite)
	}

	return errRead
}
