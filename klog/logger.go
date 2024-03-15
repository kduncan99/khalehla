// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package klog

import (
	"fmt"
	"time"
)

type Logger interface {
	Close()
	Log(level Level, message string)
	SetEnabled(bool)
	SetLevel(Level)
}

var globalEnabled = true
var globalLevel = LevelWarning
var loggers = map[Logger]*string{NewConsoleLogger(LevelWarning): nil}

func ClearLoggers() {
	loggers = make(map[Logger]*string)
}

func Close() {
	for lg, _ := range loggers {
		lg.Close()
	}
}

func RegisterLogger(logger Logger) {
	loggers[logger] = nil
}

func RegisterLoggerForSource(logger Logger, source string) {
	loggers[logger] = &source
}

func SetGlobalEnabled(enabled bool) {
	globalEnabled = enabled
}

func SetGlobalLevel(level Level) {
	globalLevel = level
}

func Log(level Level, source string, message string) {
	if level <= globalLevel && globalEnabled {
		now := time.Now()
		msg := fmt.Sprintf("%04v%02v%02v-%02v%02v%02v:%s:%s",
			now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second(), source, message)

		for lg, lgSrc := range loggers {
			if lgSrc == nil || *lgSrc == source {
				lg.Log(level, msg)
			}
		}
	}
}

func LogF(level Level, source string, format string, parameters ...interface{}) {
	if level <= globalLevel && globalEnabled {
		now := time.Now()
		msg := fmt.Sprintf("%04v%02v%02v-%02v%02v%02v:%s:",
			now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second(), source)
		msg += fmt.Sprintf(format, parameters...)

		for lg, lgSrc := range loggers {
			if lgSrc == nil || *lgSrc == source {
				lg.Log(level, msg)
			}
		}
	}
}

func LogFatal(source string, format string) {
	LogF(LevelFatal, source, format)
}

func LogFatalF(source string, format string, parameters ...interface{}) {
	LogF(LevelFatal, source, format, parameters...)
}

func LogError(source string, format string) {
	LogF(LevelError, source, format, nil)
}

func LogErrorF(source string, format string, parameters ...interface{}) {
	LogF(LevelError, source, format, parameters...)
}

func LogWarning(source string, format string) {
	LogF(LevelWarning, source, format)
}

func LogWarningF(source string, format string, parameters ...interface{}) {
	LogF(LevelWarning, source, format, parameters...)
}

func LogInfo(source string, format string) {
	LogF(LevelInfo, source, format)
}

func LogInfoF(source string, format string, parameters ...interface{}) {
	LogF(LevelInfo, source, format, parameters...)
}

func LogDebug(source string, format string) {
	LogF(LevelDebug, source, format)
}

func LogDebugF(source string, format string, parameters ...interface{}) {
	LogF(LevelDebug, source, format, parameters...)
}

func LogTrace(source string, format string) {
	LogF(LevelTrace, source, format)
}

func LogTraceF(source string, format string, parameters ...interface{}) {
	LogF(LevelTrace, source, format, parameters...)
}
