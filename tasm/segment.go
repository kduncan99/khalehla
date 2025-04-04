// khalehla Project
// tiny assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

type Segment struct {
	currentLength uint64
	generatedCode []*CodeBlock
	references    []*Reference
	labels        map[string]uint64
}

func NewSegment() *Segment {
	return &Segment{
		currentLength: 0,
		generatedCode: make([]*CodeBlock, 0),
		references:    make([]*Reference, 0),
		labels:        make(map[string]uint64),
	}
}

func (s *Segment) AppendCodeBlock(cb *CodeBlock) {
	s.generatedCode = append(s.generatedCode, cb)
	s.currentLength += uint64(len(cb.code))
}
