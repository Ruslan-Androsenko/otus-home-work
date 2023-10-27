package main

import (
	"bytes"
	"log"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	newLine := []byte("\n")
	terminalZero := []byte("\x00")

	// Читаем все файлы из каталога
	files, errDir := os.ReadDir(dir)
	if errDir != nil {
		return nil, errDir
	}

	// Инициализируем словарь для переменных среды
	environments := make(Environment, len(files))

	for _, file := range files {
		envName := file.Name()
		if file.IsDir() || strings.Contains(envName, "=") {
			continue
		}

		info, errInfo := file.Info()
		if errInfo != nil {
			log.Printf("Could not get file: %s information, error: %v \n", envName, errInfo)
			continue
		}

		// Если файл пустой, то удаляем переменную среды
		if info.Size() == 0 {
			environments[envName] = EnvValue{NeedRemove: true}
			removeEnvVariable(envName)
			continue
		}

		// Читаем содержимое из файла
		data, errData := os.ReadFile(dir + "/" + envName)
		if errData != nil {
			log.Printf("Failed to read file: %s contents, error: %v \n", envName, errData)
			continue
		}

		// Разделяем прочитанные данные на массив строк
		lines := bytes.Split(data, newLine)

		// Получаем первую строку из файла, и заменяем в ней терминальные нули на перенос строки
		firstLine := bytes.ReplaceAll(lines[0], terminalZero, newLine)

		// Очищаем полученную строку от пробелов и табуляции в конце
		envValue := strings.TrimRight(string(firstLine), " \t")

		// Если такая переменная среды уже установлена, то удаляем её
		_, ok := os.LookupEnv(envName)
		if ok {
			removeEnvVariable(envName)
		}

		// Добавляем новую переменную среды
		environments[envName] = EnvValue{Value: envValue}
	}

	return environments, nil
}

// Удалить переменную среды.
func removeEnvVariable(key string) {
	err := os.Unsetenv(key)
	if err != nil {
		log.Printf("Failed to delete environment variable: %s, error: %v \n", key, err)
	}
}
