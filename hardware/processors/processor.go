// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package processors

type ProcessorType int
type UpiIndex int

const (
	SystemProcessorType ProcessorType = iota
	InstructionProcessorType
	InputOutputProcessorType
)

type Processor interface {
	GetIndex() UpiIndex
	GetName() string
	GetType() ProcessorType
	HandleInterrupt(source UpiIndex, details interface{}) error
	Reset() error
	Start() error
	Stop()
}
