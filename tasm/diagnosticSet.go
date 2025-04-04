// khalehla Project
// tiny assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import "fmt"

type Diagnostic interface {
	GetLineNumber() uint64
	GetString() string
}

type ErrorDiagnostic struct {
	sourceSet  *SourceSet
	lineNumber uint64
	message    string
}

type InfoDiagnostic struct {
	sourceSet  *SourceSet
	lineNumber uint64
	message    string
}

type WarningDiagnostic struct {
	sourceSet  *SourceSet
	lineNumber uint64
	message    string
}

func (d *ErrorDiagnostic) GetLineNumber() uint64 {
	return d.lineNumber
}

func (d *ErrorDiagnostic) GetString() string {
	return fmt.Sprintf("E:%v:%v:%v", d.sourceSet.name, d.lineNumber, d.message)
}

func (d *InfoDiagnostic) GetLineNumber() uint64 {
	return d.lineNumber
}

func (d *InfoDiagnostic) GetString() string {
	return fmt.Sprintf("I:%v:%v:%v", d.sourceSet.name, d.lineNumber, d.message)
}

func (d *WarningDiagnostic) GetLineNumber() uint64 {
	return d.lineNumber
}

func (d *WarningDiagnostic) GetString() string {
	return fmt.Sprintf("W:%v:%v:%v", d.sourceSet.name, d.lineNumber, d.message)
}

type DiagnosticSet struct {
	diagnostics  map[uint64][]Diagnostic
	infoCount    uint64
	warningCount uint64
	errorCount   uint64
}

func NewDiagnosticSet() *DiagnosticSet {
	return &DiagnosticSet{
		diagnostics: make(map[uint64][]Diagnostic),
	}
}

func (ds *DiagnosticSet) putDiag(diag Diagnostic) {
	slice, ok := ds.diagnostics[diag.GetLineNumber()]
	if ok {
		ds.diagnostics[diag.GetLineNumber()] = append(slice, diag)
	} else {
		ds.diagnostics[diag.GetLineNumber()] = []Diagnostic{diag}
	}
}

func (ds *DiagnosticSet) NewError(source *SourceSet, lineNumber uint64, message string) {
	d := &ErrorDiagnostic{
		sourceSet:  source,
		lineNumber: lineNumber,
		message:    message,
	}
	ds.putDiag(d)
	ds.errorCount++
}

func (ds *DiagnosticSet) NewInfo(source *SourceSet, lineNumber uint64, message string) {
	d := &InfoDiagnostic{
		sourceSet:  source,
		lineNumber: lineNumber,
		message:    message,
	}

	ds.putDiag(d)
	ds.infoCount++
}

func (ds *DiagnosticSet) NewWarning(source *SourceSet, lineNumber uint64, message string) {
	d := &WarningDiagnostic{
		sourceSet:  source,
		lineNumber: lineNumber,
		message:    message,
	}

	ds.putDiag(d)
	ds.infoCount++
}
