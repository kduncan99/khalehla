// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package klog

import (
	"fmt"
	"io"
	"os"
	"time"
)

type TimestampedFileLogger struct {
	enabled bool
	level   Level
	file    *os.File
	writer  io.Writer
}

func NewTimestampedFileLogger(initialLevel Level, prefix string) *TimestampedFileLogger {
	now := time.Now()
	filename := fmt.Sprintf("%v-%04v%02v%02v-%02v%02v%02v.log",
		prefix, now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0666)
	if err != nil {
		panic(err.Error())
	}

	return &TimestampedFileLogger{
		enabled: true,
		level:   initialLevel,
		file:    file,
	}
}

func (lg *TimestampedFileLogger) Close() {
	_ = lg.file.Close()
}

func (lg *TimestampedFileLogger) Log(level Level, message string) {
	if level <= lg.level && lg.enabled {
		_, _ = lg.file.WriteString(message + "\n")
	}
}

func (lg *TimestampedFileLogger) SetEnabled(enabled bool) {
	lg.enabled = enabled
}

func (lg *TimestampedFileLogger) SetLevel(level Level) {
	lg.level = level
}
