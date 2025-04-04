// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tapeUtil

type FileHeader struct {
	// TODO so far, nothing we handle produces or creates a file header on media,
	// so we don't actually know what this struct might contain.
}

func NewFileHeader() *FileHeader {
	return &FileHeader{}
}
