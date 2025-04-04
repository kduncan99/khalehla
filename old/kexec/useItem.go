// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

// UseItem for keeping track of @USE for a run.
// Ensure that the fileSpec has a qualifier...
// i.e., resolve implied or default qualifiers before posting the useItem.
type UseItem struct {
	InternalFilename  string
	FileSpecification *FileSpecification
	ReleaseFlag       bool
}
