// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package processors

import (
	"fmt"
	"sync"

	"khalehla/logger"
)

// SystemProcessor manages all the other processors. There is only one SystemProcessor in any configuration.
type SystemProcessor struct {
	upiIndex   UpiIndex
	name       string
	processors map[UpiIndex]Processor // map of all Processor entities (including ourself)
	mutex      sync.Mutex
}

func NewSystemProcessor() *SystemProcessor {
	p := new(SystemProcessor)
	p.name = "SP0"
	p.upiIndex = 0
	p.processors = make(map[UpiIndex]Processor)
	p.processors[p.upiIndex] = p
	return p
}

func (sp *SystemProcessor) GetIndex() UpiIndex {
	return sp.upiIndex
}

func (sp *SystemProcessor) GetName() string {
	return sp.name
}

func (sp *SystemProcessor) GetType() ProcessorType {
	return SystemProcessorType
}

// GetProcessor retrieves the processor for the given UPI Index.
// Mainly for use by other processors.
func (sp *SystemProcessor) GetProcessor(upiIndex UpiIndex) (Processor, error) {
	proc, ok := sp.processors[upiIndex]
	if !ok {
		return nil, fmt.Errorf("processor for upi %v not found", upiIndex)
	}
	return proc, nil
}

// HandleInterrupt handles any UPI sent to us from some other processor.
// We only accept interrupts from InstructionProcessor entities, which indicate to us that the IP has halted.
func (sp *SystemProcessor) HandleInterrupt(source UpiIndex, details interface{}) error {
	proc, err := sp.GetProcessor(source)
	if err != nil {
		return err
	}
	switch proc.GetType() {
	case InstructionProcessorType:
		// TODO (halted)
	default:
		msg := fmt.Sprintf("interrupt from source %v not handled", source)
		logger.LogFatal(sp.name, msg)
		return fmt.Errorf(msg)
	}

	return nil
}

// SendInterrupt is intended to be invoked by IP and IOP (and by ourselves) to ping some other Processor.
//
//	In practice, the following messages are thus represented:
//		SYS -> IP starts the instruction processor - details indicate a particular HardwareInterrupt to be invoked
//		IOP -> IP indicates that an IO operation has completed - details is the AbsoluteAddress of the corresponding ChannelProgram
//		IP -> IOP indicates that an IO operation is to be started - details is the AbsoluteAddress of the corresponding ChannelProgram
//		IP -> SYS indicates that an IP has halted - details indicates the reason why
func (sp *SystemProcessor) SendInterrupt(source UpiIndex, destination UpiIndex, details interface{}) error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()
	processor, ok := sp.processors[source]
	if !ok {
		return fmt.Errorf("destination processor %s not found for source %v", destination, source)
	}

	return processor.HandleInterrupt(source, details)
}

func (sp *SystemProcessor) Reset() (err error) {
	logger.Log(logger.LevelTrace, sp.name, "Reset")
	return
}

func (sp *SystemProcessor) Start() (err error) {
	logger.Log(logger.LevelTrace, sp.name, "Start")
	return
}

func (sp *SystemProcessor) Stop() {
	logger.Log(logger.LevelTrace, sp.name, "Stop")
	return
}

// CreateStorageComplex builds up the tree of processors given the number of desired IPs and IOPs
func (sp *SystemProcessor) CreateStorageComplex(ipCount int, iopCount int) error {
	if ipCount < 1 || iopCount < 1 {
		return fmt.Errorf("ip and iop count must be greater than 0")
	}

	for _, proc := range sp.processors {
		proc.Stop()
		err := proc.Reset()
		if err != nil {
			logger.LogWarningF(sp.name, "Failed to reset %v while removing it: %v", proc.GetName(), err)
		}
	}
	sp.processors = make(map[UpiIndex]Processor)

	upix := sp.upiIndex + 1
	for ix := 0; ix < ipCount; ix++ {
		ip := NewInstructionProcessor(upix, fmt.Sprintf("IP%v", ix), sp)
		sp.processors[ip.upiIndex] = ip
		upix++
	}

	for ix := 0; ix < iopCount; ix++ {
		iop := NewInputOutputProcessor(upix, fmt.Sprintf("IOP%v", ix), sp)
		sp.processors[iop.upiIndex] = iop
		upix++
	}

	return nil
}

// TODO need means of adding channels/devices to IOPs
