package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if len(fromPath) == 0 {
		return errors.New("from path a file does not exist")
	}

	if len(toPath) == 0 {
		return errors.New("to path a file does not exist")
	}

	input, errInput := os.Open(fromPath)
	if errInput != nil {
		errMessage := fmt.Sprintf("failed read from input file, error: %v", errInput)
		return errors.New(errMessage)
	}

	info, errInfo := input.Stat()
	if errInfo != nil {
		errMessage := fmt.Sprintf("failed read info meta data from input file, error: %v", errInfo)
		return errors.New(errMessage)
	}

	inputFileSize := info.Size()
	if inputFileSize <= 0 {
		return ErrUnsupportedFile
	}

	if offset > inputFileSize {
		return ErrOffsetExceedsFileSize
	}

	output, errOutput := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if errOutput != nil {
		errMessage := fmt.Sprintf("failed write to output file, error: %v", errOutput)
		return errors.New(errMessage)
	}

	var (
		hasEndWrite bool
		writeOffset int64
	)

	bufferSize := 1024
	buffer := make([]byte, bufferSize)
	clearBuffer := make([]byte, bufferSize)

	for offset < inputFileSize {
		read, errRead := input.ReadAt(buffer, offset)
		if errRead != nil && errRead != io.EOF {
			errMessage := fmt.Sprintf("error reading from input file, error: %v", errRead)
			return errors.New(errMessage)
		}

		// Если прочитали меньше чем ожидалось, то уменьшаем размер буфера
		if read < bufferSize {
			buffer = buffer[:read]
			hasEndWrite = true
		} else if limit > 0 && limit < int64(read) {
			// Если заданый лимит меньше чем прочитаная часть данных, то уменьшаем размер буфера
			buffer = buffer[:limit]
			hasEndWrite = true
		} else if limit > 0 && limit < writeOffset+int64(bufferSize) {
			// Если текущий проход записи данных, превышает заданый лимит, то уменьшаем размер буфера
			sliceOffset := int64(math.Abs(float64(writeOffset - limit)))
			buffer = buffer[:sliceOffset]
			hasEndWrite = true
		}

		_, errWrite := output.WriteAt(buffer, writeOffset)
		if errWrite != nil {
			errMessage := fmt.Sprintf("error writing to output file, error: %v", errWrite)
			return errors.New(errMessage)
		}

		// Очищаем буфер с данными
		copy(buffer, clearBuffer)

		// Смещаем позицию в файле
		offset += int64(read)
		writeOffset += int64(read)

		// Если файл дочитан до конца, или установлен флаг
		// что запрошенный объем данных уже был записан, то выходим из цикла
		if errRead == io.EOF || hasEndWrite {
			break
		}
	}

	errInput = input.Close()
	if errInput != nil {
		return errInput
	}

	errOutput = output.Close()
	if errOutput != nil {
		return errOutput
	}

	return nil
}
