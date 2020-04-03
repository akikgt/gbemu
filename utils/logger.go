package utils

import "fmt"

type Logger struct {
	enabled bool
}

func NewLogger(enabled bool) *Logger {
	logger := &Logger{enabled}
	return logger
}

func (logger *Logger) Log(format string, params ...interface{}) {
	if !logger.enabled {
		return
	}

	fmt.Printf(format, params...)
}
