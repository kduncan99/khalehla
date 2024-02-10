// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"sync"
	"time"
)

type LogLevel int

const (
	LogFatal   LogLevel = 0
	LogError   LogLevel = 1
	LogWarning LogLevel = 2
	LogInfo    LogLevel = 3
	LogDebug   LogLevel = 9
	LogTrace   LogLevel = 10
)

type LogDestination interface {
	CreateEntry(level LogLevel, source string, message string)
}

type LogEntry struct {
	time    time.Time
	level   LogLevel
	source  string
	message string
}

type Logger struct {
	mutex        sync.Mutex
	logLevel     LogLevel
	destinations []LogDestination
}

var SystemLogger = NewLogger()

func NewLogger() *Logger {
	logger := Logger{}
	logger.destinations = make([]LogDestination, 0)
	return &logger
}

func (l *Logger) CreateEntry(level LogLevel, source string, message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if level <= l.logLevel {
		for _, dest := range l.destinations {
			dest.CreateEntry(level, source, message)
		}
	}
}

func (l *Logger) CreateEntryAlways(level LogLevel, source string, message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for _, dest := range l.destinations {
		dest.CreateEntry(level, source, message)
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.logLevel = level
}

func (l *Logger) AddDestination(dest LogDestination) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.destinations = append(l.destinations, dest)
}
