// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package klog

import "fmt"

type ConsoleLogger struct {
	enabled bool
	level   Level
}

func NewConsoleLogger(initialLevel Level) *ConsoleLogger {
	return &ConsoleLogger{
		enabled: true,
		level:   initialLevel,
	}
}

func (lg *ConsoleLogger) Close() {
	// nothing to do
}

func (lg *ConsoleLogger) Log(level Level, message string) {
	if level <= lg.level && lg.enabled {
		fmt.Println(message)
	}
}

func (lg *ConsoleLogger) SetEnabled(enabled bool) {
	lg.enabled = enabled
}

func (lg *ConsoleLogger) SetLevel(level Level) {
	lg.level = level
}
