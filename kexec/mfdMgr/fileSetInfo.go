// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

type FileSetCycleInfo struct {
	toBeCataloged bool
	toBeDropped   bool
	absoluteCycle uint
}

type FileSetInfo struct {
	qualifier       string
	filename        string
	projectId       string
	readKey         string
	writeKey        string
	fileType        FileType
	plusOneExists   bool
	count           uint
	maxCycleRange   uint
	currentRange    uint
	highestAbsolute uint
	fileInfo        []FileSetCycleInfo
}
