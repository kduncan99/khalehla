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

func (lg *ConsoleLogger) Log(level Level, message string) {
	if level <= globalLevel && globalEnabled {
		fmt.Println(message)
	}
}

func (lg *ConsoleLogger) SetEnabled(enabled bool) {
	lg.enabled = enabled
}

func (lg *ConsoleLogger) SetLevel(level Level) {
	lg.level = level
}
