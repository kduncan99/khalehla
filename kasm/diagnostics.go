// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type diagnosticType int

const (
	INFO diagnosticType = iota
	WARNING
	ERROR
)

type diagnosticEntry struct {
	diagType diagnosticType
	text     string
}

func (d *diagnosticEntry) GetDiagnosticType() diagnosticType {
	return d.diagType
}

func (d *diagnosticEntry) GetText() string {
	return d.text
}

type diagnostics struct {
	// key is line number, value is array of diagnostics pertaining to that line number
	entries map[int][]diagnosticEntry
}

func (d *diagnostics) Clear() {
	d.entries = make(map[int][]diagnosticEntry)
}

func (d *diagnostics) AppendDiagnostic(lineNumber int, entry diagnosticEntry) {
	_, ok := d.entries[lineNumber]
	if !ok {
		d.entries[lineNumber] = []diagnosticEntry{entry}
	} else {
		d.entries[lineNumber] = append(d.entries[lineNumber], entry)
	}
}

func (d *diagnostics) AppendError(lineNumber int, text string) {
	diag := diagnosticEntry{
		diagType: ERROR,
		text:     text,
	}
	d.AppendDiagnostic(lineNumber, diag)
}

func (d *diagnostics) AppendInfo(lineNumber int, text string) {
	diag := diagnosticEntry{
		diagType: INFO,
		text:     text,
	}
	d.AppendDiagnostic(lineNumber, diag)
}

func (d *diagnostics) AppendWarning(lineNumber int, text string) {
	diag := diagnosticEntry{
		diagType: WARNING,
		text:     text,
	}
	d.AppendDiagnostic(lineNumber, diag)
}

func (d *diagnostics) GetDiagnosticCounters() (int, int, int) {
	errors := 0
	infos := 0
	warnings := 0

	for _, arr := range d.entries {
		for _, diag := range arr {
			if diag.GetDiagnosticType() == ERROR {
				errors++
			} else if diag.GetDiagnosticType() == INFO {
				infos++
			} else if diag.GetDiagnosticType() == WARNING {
				warnings++
			}
		}
	}

	return errors, warnings, infos
}
