// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// BackupInfo describes information from the file cycle pertaining to the backup state of the file.
// It is contained within one or more of the more-specific file cycle info structs.
// We do not (currently) support multiple backup levels, so don't look for that information.
type BackupInfo struct {
	TimeBackupCreated    uint64
	FASBits              uint64
	NumberOfTextBlocks   uint64
	StartingFilePosition uint64
	BackupReelNumbers    []string
}
