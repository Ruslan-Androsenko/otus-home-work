package logger

import (
	"fmt"
	"os"
	"time"
)

type (
	logLevel uint8
	Logger   struct {
		level logLevel
	}
)

func New(level string) *Logger {
	logg := Logger{
		level: errorLevel,
	}

	if levelIndex, ok := titleToLevel[level]; ok {
		logg.level = levelIndex
	}

	return &logg
}

func (l Logger) Fatal(msg string) {
	l.printMessage(msg, fatalLevel)
	os.Exit(1)
}

func (l Logger) Error(msg string) {
	l.printMessage(msg, errorLevel)
}

func (l Logger) Warning(msg string) {
	l.printMessage(msg, warningLevel)
}

func (l Logger) Info(msg string) {
	l.printMessage(msg, infoLevel)
}

func (l Logger) Debug(msg string) {
	l.printMessage(msg, debugLevel)
}

// Напечатать отформатированное сообщение в консоль.
func (l Logger) printMessage(msg string, level logLevel) {
	if level > l.level {
		return
	}

	logTime := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Printf("[%s] [%s] %s \n", logTime, level, msg)
}
