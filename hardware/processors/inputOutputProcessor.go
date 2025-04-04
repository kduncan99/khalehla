// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package processors

import (
	"fmt"

	"khalehla/logger"
	"khalehla/old/hardware/channels"
)

// An InputOutputProcessor responds to UPI messages from an InstructionProcessor.
// Such a message will be accompanied by an AbsoluteAddress indicating the location in memory of a
// ChannelProgram which is passed to a particular channel for processing.
// Our job is to notify the indicated channel of the existence of the channel program, and to monitor the
// channel program until it is complete, whereupon we send a UPI message to any available InstructionProcessor
// which completes the IO operation.
type InputOutputProcessor struct {
	sp       *SystemProcessor
	upiIndex UpiIndex
	name     string
	channels map[int]channels.Channel
}

func NewInputOutputProcessor(index UpiIndex, name string, systemProcessor *SystemProcessor) *InputOutputProcessor {
	p := new(InputOutputProcessor)
	p.sp = systemProcessor
	p.name = name
	p.upiIndex = index
	return p
}

func (iop *InputOutputProcessor) GetIndex() UpiIndex {
	return iop.upiIndex
}

func (iop *InputOutputProcessor) GetName() string {
	return iop.name
}

func (iop *InputOutputProcessor) GetType() ProcessorType {
	return InputOutputProcessorType
}

func (iop *InputOutputProcessor) HandleInterrupt(source UpiIndex, details interface{}) error {
	proc, err := iop.sp.GetProcessor(source)
	if err != nil {
		return err
	}
	switch proc.GetType() {
	case InstructionProcessorType:
		// TODO (IO complete)
	default:
		msg := fmt.Sprintf("interrupt from source %v not handled", source)
		logger.LogFatal(iop.name, msg)
		return fmt.Errorf(msg)
	}

	return nil
}

func (iop *InputOutputProcessor) Reset() (err error) {
	// TODO
	return
}

func (iop *InputOutputProcessor) Start() (err error) {
	// TODO
	return
}

func (iop *InputOutputProcessor) Stop() {
	// TODO
}

// TODO function to adopt channels, manage channels, send IO to channels, and to manage UPI messages
