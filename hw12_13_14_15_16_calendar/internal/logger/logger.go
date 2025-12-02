package logger

import "fmt"

type Logger struct {
	level    string
	location string
}

func New(level string, location string) *Logger {
	return &Logger{level: level,
		location: location}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	fmt.Println(msg)
}
