// Khalehla Project
// testing assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import "khalehla/pkg"

// Module is a single assembled module, from one assembly.
// It must be bundled into a Loadable before being loaded and presented to the emulator.
type Module struct {
	code map[int][]pkg.Word36
}

func (m *Module) assemble(sourceName string, source []string) (result bool, diagnostics []string) {
	result = true
	diagnostics = make([]string, 0)

	parsedLines := make([]*parsedLine, len(source))
	for sx := 0; sx < len(source); sx++ {
		parsedLines[sx] = parse(&sourceName, sx+1, source[sx])
		if !parsedLines[sx].result {
			result = false
		}
		diagnostics = append(diagnostics, parsedLines[sx].diagnostics...)
	}

	return
}
