// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package processors

import (
	"fmt"

	"khalehla/hardware/processors/ipEngine"
	"khalehla/logger"
)

// An InstructionProcessor executes 36-bit architecturally-defined code.
type InstructionProcessor struct {
	sp       *SystemProcessor
	upiIndex UpiIndex
	name     string
	engine   *ipEngine.InstructionEngine
}

func NewInstructionProcessor(index UpiIndex, name string, systemProcessor *SystemProcessor) *InstructionProcessor {
	p := new(InstructionProcessor)
	p.sp = systemProcessor
	p.name = name
	p.upiIndex = index
	return p
}

func (ip *InstructionProcessor) GetIndex() UpiIndex {
	return ip.upiIndex
}

func (ip *InstructionProcessor) GetName() string {
	return ip.name
}

func (ip *InstructionProcessor) GetType() ProcessorType {
	return InstructionProcessorType
}

func (ip *InstructionProcessor) HandleInterrupt(source UpiIndex, details interface{}) error {
	proc, err := ip.sp.GetProcessor(source)
	if err != nil {
		return err
	}
	switch proc.GetType() {
	case InputOutputProcessorType:
		// TODO (IO complete)
	case SystemProcessorType:
		// TODO (Initial start, start, or stop) - do we need this?
	default:
		msg := fmt.Sprintf("interrupt from source %v not handled", source)
		logger.LogFatal(ip.name, msg)
		return fmt.Errorf(msg)
	}

	return nil
}

func (ip *InstructionProcessor) Reset() (err error) {
	// TODO
	return
}

func (ip *InstructionProcessor) Start() (err error) {
	// TODO
	return
}

func (ip *InstructionProcessor) Stop() {
	// TODO
}

// TODO function to adopt channels, manage channels, send IO to channels, and to manage UPI messages
