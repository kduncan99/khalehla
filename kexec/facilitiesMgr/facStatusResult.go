// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import "khalehla/kexec"

type FacStatusMessageInstance struct {
	code   kexec.FacStatusCode
	values []string
}

type FacStatusResult struct {
	Infos    []*FacStatusMessageInstance
	Warnings []*FacStatusMessageInstance
	Errors   []*FacStatusMessageInstance
}

func NewFacResult() *FacStatusResult {
	return &FacStatusResult{
		Infos:    make([]*FacStatusMessageInstance, 0),
		Warnings: make([]*FacStatusMessageInstance, 0),
		Errors:   make([]*FacStatusMessageInstance, 0),
	}
}

func (fr *FacStatusResult) HasInformationalMessages() bool {
	return len(fr.Infos) > 0
}

func (fr *FacStatusResult) HasWarningMessages() bool {
	return len(fr.Warnings) > 0
}

func (fr *FacStatusResult) HasErrorMessages() bool {
	return len(fr.Errors) > 0
}

func (fr *FacStatusResult) PostMessage(code kexec.FacStatusCode, values []string) {
	msg := &FacStatusMessageInstance{
		code:   code,
		values: values,
	}
	temp := FacStatusMessageTemplates[code]
	switch temp.Category {
	case kexec.FacMsgInfo:
		fr.Infos = append(fr.Infos, msg)
	case kexec.FacMsgWarning:
		fr.Warnings = append(fr.Warnings, msg)
	case kexec.FacMsgError:
		fr.Errors = append(fr.Errors, msg)
	}
}
