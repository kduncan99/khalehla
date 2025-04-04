// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package logger

type Level uint

const (
	LevelSilent Level = iota
	LevelFatal
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
	LevelAll Level = 99
)

var LevelLookup = map[Level]string{
	LevelFatal:   "FATAL",
	LevelError:   "ERROR",
	LevelWarning: "WARN",
	LevelInfo:    "INFO",
	LevelDebug:   "DEBUG",
	LevelTrace:   "TRACE",
}
