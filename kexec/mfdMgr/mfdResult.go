// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

type MFDResult uint

const (
	MFDSuccessful MFDResult = iota
	MFDInternalError
	MFDFileNameConflict
	MFDNotFound
)
