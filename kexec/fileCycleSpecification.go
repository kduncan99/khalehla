// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type FileCycleSpecification struct {
	AbsoluteCycle *uint
	RelativeCycle *int
}

func NewAbsoluteFileCycleSpecification(cycle uint) *FileCycleSpecification {
	ac := cycle
	return &FileCycleSpecification{
		AbsoluteCycle: &ac,
		RelativeCycle: nil,
	}
}

func NewRelativeFileCycleSpecification(cycle int) *FileCycleSpecification {
	rc := cycle
	return &FileCycleSpecification{
		AbsoluteCycle: nil,
		RelativeCycle: &rc,
	}
}

func CopyFileCycleSpecification(fcs *FileCycleSpecification) *FileCycleSpecification {
	return &FileCycleSpecification{
		AbsoluteCycle: fcs.AbsoluteCycle,
		RelativeCycle: fcs.RelativeCycle,
	}
}

func (fcs *FileCycleSpecification) IsRelative() bool {
	return fcs.RelativeCycle != nil
}

func (fcs *FileCycleSpecification) IsAbsolute() bool {
	return fcs.AbsoluteCycle != nil
}
