// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package klog

import (
	"fmt"
	"time"
)

type Logger interface {
	Log(level Level, message string)
	SetEnabled(bool)
	SetLevel(Level)
}

var globalEnabled = true
var globalLevel = LevelWarning
var loggers = map[Logger]*string{NewConsoleLogger(LevelWarning): nil}

func ClearLogger() {
	loggers = make(map[Logger]*string)
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

func Log(level Level, source string, format string, parameters ...interface{}) {
	if level <= globalLevel && globalEnabled {
		now := time.Now()
		msg := fmt.Sprintf("%04v%02v%02v-%02v%02v%02v:%s:",
			now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second(), source)
		if parameters != nil {
			msg += fmt.Sprintf(format, parameters...)
		} else {
			msg += format
		}

		for lg, lgSrc := range loggers {
			if lgSrc == nil || *lgSrc == source {
				lg.Log(level, msg)
			}
		}
	}
}

func LogFatal(source string, format string) {
	Log(LevelFatal, source, format, nil)
}

func LogFatalF(source string, format string, parameters ...interface{}) {
	Log(LevelFatal, source, format, parameters...)
}

func LogError(source string, format string) {
	Log(LevelError, source, format, nil)
}

func LogErrorF(source string, format string, parameters ...interface{}) {
	Log(LevelError, source, format, parameters...)
}

func LogWarning(source string, format string) {
	Log(LevelWarning, source, format, nil)
}

func LogWarningF(source string, format string, parameters ...interface{}) {
	Log(LevelWarning, source, format, parameters...)
}

func LogInfo(source string, format string) {
	Log(LevelInfo, source, format, nil)
}

func LogInfoF(source string, format string, parameters ...interface{}) {
	Log(LevelInfo, source, format, parameters...)
}

func LogDebug(source string, format string) {
	Log(LevelDebug, source, format, nil)
}

func LogDebugF(source string, format string, parameters ...interface{}) {
	Log(LevelDebug, source, format, parameters...)
}

func LogTrace(source string, format string) {
	Log(LevelTrace, source, format, nil)
}

func LogTraceF(source string, format string, parameters ...interface{}) {
	Log(LevelTrace, source, format, parameters...)
}
