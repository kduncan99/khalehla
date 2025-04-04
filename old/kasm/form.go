// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import "fmt"

type Form struct {
	bitSizes []int
}

var SimpleForm, _ = NewForm([]int{36})

func sumOf(values []int) int {
	sum := 0
	for _, i := range values {
		sum += i
	}
	return sum
}

func NewForm(bitSizes []int) (*Form, error) {
	if sumOf(bitSizes) == 36 {
		f := &Form{
			bitSizes: bitSizes,
		}

		return f, nil
	}

	return nil, fmt.Errorf("incomplete form")
}

func (f *Form) Equals(op *Form) bool {
	if len(f.bitSizes) == len(op.bitSizes) {
		for bx := 0; bx < len(f.bitSizes); bx++ {
			if f.bitSizes[bx] != op.bitSizes[bx] {
				return false
			}
		}

		return true
	}

	return false
}
