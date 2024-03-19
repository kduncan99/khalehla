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
		msg := fmt.Sprintf("%04v%02v%02v-%02v%02v%02v:%s:%s:%s",
			now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second(),
			LevelLookup[level], source, message)

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
		msg := fmt.Sprintf("%04v%02v%02v-%02v%02v%02v:%s:%s:",
			now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second(),
			LevelLookup[level], source)
		msg += fmt.Sprintf(format, parameters...)

		for lg, lgSrc := range loggers {
			if lgSrc == nil || *lgSrc == source {
				lg.Log(level, msg)
			}
		}
	}
}

func LogFatal(source string, message string) {
	Log(LevelFatal, source, message)
}

func LogFatalF(source string, format string, parameters ...interface{}) {
	LogF(LevelFatal, source, format, parameters...)
}

func LogError(source string, message string) {
	Log(LevelError, source, message)
}

func LogErrorF(source string, format string, parameters ...interface{}) {
	LogF(LevelError, source, format, parameters...)
}

func LogWarning(source string, message string) {
	Log(LevelWarning, source, message)
}

func LogWarningF(source string, format string, parameters ...interface{}) {
	LogF(LevelWarning, source, format, parameters...)
}

func LogInfo(source string, message string) {
	Log(LevelInfo, source, message)
}

func LogInfoF(source string, format string, parameters ...interface{}) {
	LogF(LevelInfo, source, format, parameters...)
}

func LogDebug(source string, message string) {
	Log(LevelDebug, source, message)
}

func LogDebugF(source string, format string, parameters ...interface{}) {
	LogF(LevelDebug, source, format, parameters...)
}

func LogTrace(source string, message string) {
	Log(LevelTrace, source, message)
}

func LogTraceF(source string, format string, parameters ...interface{}) {
	LogF(LevelTrace, source, format, parameters...)
}
