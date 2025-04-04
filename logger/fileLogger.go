// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package logger

import (
	"os"
)

type FileLogger struct {
	enabled bool
	level   Level
	file    *os.File
}

func NewFileLogger(initialLevel Level, filename string) *FileLogger {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err.Error())
	}

	return &FileLogger{
		enabled: true,
		level:   initialLevel,
		file:    file,
	}
}

func (lg *FileLogger) Close() {
	_ = lg.file.Close()
}

func (lg *FileLogger) Log(level Level, message string) {
	if level <= lg.level && lg.enabled {
		_, _ = lg.file.Write([]byte(message + "/n"))
	}
}

func (lg *FileLogger) SetEnabled(enabled bool) {
	lg.enabled = enabled
}

func (lg *FileLogger) SetLevel(level Level) {
	lg.level = level
}
