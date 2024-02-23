// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

type BackupInfo struct {
	timeBackupCreated          uint64
	maxBackupLevels            uint64
	currentBackupLevels        uint64
	fasBits                    uint64
	numberOfTextBlocks         uint64
	backupStartingFilePosition uint64
	backupReelNumbers          []string
}
