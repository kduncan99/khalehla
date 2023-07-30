// Khalehla Project
// tiny assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import (
	"fmt"
	"strings"
)

type SourceSet struct {
	name        string
	sourceItems []*SourceItem
}

type SourceItem struct {
	label    *string
	command  *string
	operands []string
}

func NewSourceItem(label string, command string, operands []string) *SourceItem {
	upLabel := strings.ToUpper(label)
	upCommand := strings.ToUpper(command)

	return &SourceItem{
		label:    &upLabel,
		command:  &upCommand,
		operands: operands,
	}
}

func (si *SourceItem) GetString() string {
	lab := ""
	if si.label != nil {
		lab = *si.label
	}

	cmd := ""
	if si.command != nil {
		cmd = *si.command
	}

	ops := ""
	if si.operands != nil {
		ops = strings.Join(si.operands, ",")
	}

	return fmt.Sprintf("%-10s %-10s %s", lab, cmd, ops)
}

func NewSourceSet(name string, sourceItems []*SourceItem) *SourceSet {
	return &SourceSet{
		name:        name,
		sourceItems: sourceItems,
	}
}
