// Khalehla Project
// tiny assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

type Segment struct {
	currentLength int
	generatedCode []*CodeBlock
	references    []*Reference
	labels        map[string]int
}

func NewSegment() *Segment {
	return &Segment{
		currentLength: 0,
		generatedCode: make([]*CodeBlock, 0),
		references:    make([]*Reference, 0),
		labels:        make(map[string]int),
	}
}

func (s *Segment) AppendCodeBlock(cb *CodeBlock) {
	s.generatedCode = append(s.generatedCode, cb)
	s.currentLength += len(cb.code)
}
