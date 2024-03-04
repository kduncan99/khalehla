// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package config

import (
	"fmt"
)

type EquipmentUsage uint

const (
	EquipmentUsageSectorAddressableMassStorage = iota
	EquipmentUsageWordAddressableMassStorage
	EquipmentUsageTape
)

type EquipmentEntry struct {
	Mnemonic            string
	Usage               EquipmentUsage
	SelectableEquipment []string // name associated with the NodeModelType (see nodeMgr)
}

type Configuration struct {
	AccountInitialReserve          uint64 // initial reserve for SYS$*ACCOUNT$R1 and SYS$SEC@ACCTINFO files
	AccountAssignMnemonic          string // assign mnemonic for SYS$*ACCOUNT$R1 and SYS$SEC@ACCTINFO files
	DLOCAssignMnemonic             string // assign mnemonic for SYS$*DLOC$ file
	FilesPrivateByAccount          bool
	GenFInitialReserve             uint64 // initial reserve for SYS$*GENF$ file
	GenFAssignMnemonic             string // assign mnemonic for SYS$*GENF$ file
	LibInitialReserve              uint64 // initial reserve for SYS$*LIB$ file
	LibAssignMnemonic              string // assign mnemonic for SYS$*LIB$ file
	LibMaximumSize                 uint64 // max granules for SYS$*LIB$ file
	LogConsoleMessages             bool
	LogIOs                         bool
	LogTrace                       bool
	MasterAccountId                string // could be empty, in which case operator is prompted when ACCOUNT$R1 is created
	MassStorageDefaultMnemonic     string // Usually 'F'
	MaxCards                       uint64
	MaxGranules                    uint64 // max granules if not specified on @ASG or @CAT
	MaxPages                       uint64
	OverheadAccountId              string // account ID for overhead runs such as SYS and ROLOUT/ROLBACK
	OverheadUserId                 string // User ID for overhead runs
	PrivilegedAccountId            string // account ID which can override reading tape label blocks
	ReleaseUnusedReserve           bool
	ReleaseUnusedRemovableReserve  bool
	ResidueClear                   bool   // zero out tracks when allocated
	RunInitialReserve              uint64 // initial reserve for SYS$*RUN$ file
	RunAssignMnemonic              string // assign mnemonic for SYS$*RUN$ file
	RunMaximumSize                 uint64 // max granules for SYS$*RUN$ file
	SacrdInitialReserve            uint64 // initial reserve for SYS$*SEC@ACR$ file
	SacrdAssignMnemonic            string // assign mnemonic for SYS$*SEC@ACR$ file
	SecurityOfficerUserId          string // could be empty, in which case operator is prompted at boot time
	SymbiontBufferSize             uint64 // Buffer size used for standard and alternate read/write buffers
	SystemTapeEquipment            string // assign mnemonic for exec tape requests
	TapeAccessRestrictedByAccount  bool
	TapeDefaultMnemonic            string
	TerminateMaxCards              bool
	TerminateMaxPages              bool
	TerminateMaxTime               bool
	TIPQualifier                   string
	TIPReadKey                     string
	TIPWriteKey                    string
	TPFAssignMnemonic              string
	TPFMaxSize                     uint64
	UserInitialReserve             uint64 // initial reserve for SYS$*SEC@USERID$ file
	UserAssignMnemonic             string // assign mnemonic for SYS$*SEC@USERID$ file
	WordAddressableDefaultMnemonic string
	EquipmentTable                 map[string]*EquipmentEntry // key is mnemonic

	// TODO -- and this is not exhaustive...
	// SMDTFASGMNE SYS$*SMDTF$ (?) assign mnemonic
	// SMDTFINTRES (ditto) initial reserve [1]
	// TLAUTO automatic tape labeling [false]
}

func NewConfiguration() *Configuration {
	cfg := &Configuration{}

	cfg.AccountInitialReserve = 10
	cfg.AccountAssignMnemonic = "F"
	cfg.DLOCAssignMnemonic = "F"
	cfg.FilesPrivateByAccount = true
	cfg.GenFInitialReserve = 128
	cfg.GenFAssignMnemonic = "F"
	cfg.LibInitialReserve = 128
	cfg.LibAssignMnemonic = "F"
	cfg.LibMaximumSize = 9999
	cfg.LogConsoleMessages = true
	cfg.LogIOs = true   // TODO false
	cfg.LogTrace = true // TODO false
	cfg.MassStorageDefaultMnemonic = "F"
	cfg.MasterAccountId = "SYSTEM"
	cfg.MaxCards = 256
	cfg.MaxGranules = 256
	cfg.MaxPages = 256
	cfg.OverheadUserId = "INSTALLATION"
	cfg.OverheadAccountId = "INSTALLATION"
	cfg.PrivilegedAccountId = "123456"
	cfg.ReleaseUnusedReserve = true
	cfg.ReleaseUnusedRemovableReserve = false
	cfg.ResidueClear = true
	cfg.RunInitialReserve = 10
	cfg.RunAssignMnemonic = "F"
	cfg.RunMaximumSize = 256
	cfg.SacrdInitialReserve = 10
	cfg.SacrdAssignMnemonic = "F"
	cfg.SecurityOfficerUserId = ""
	cfg.SymbiontBufferSize = 224
	cfg.SystemTapeEquipment = "T"
	cfg.TapeAccessRestrictedByAccount = false
	cfg.TapeDefaultMnemonic = "T"
	cfg.TerminateMaxCards = false
	cfg.TerminateMaxPages = false
	cfg.TerminateMaxTime = false
	cfg.TIPQualifier = "TIP$"
	cfg.TIPReadKey = "++++++"
	cfg.TIPWriteKey = "++++++"
	cfg.TPFAssignMnemonic = "F"
	cfg.TPFMaxSize = 128
	cfg.UserInitialReserve = 10
	cfg.UserAssignMnemonic = "F"
	cfg.WordAddressableDefaultMnemonic = "D"

	cfg.EquipmentTable = make(map[string]*EquipmentEntry)

	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "F",
		Usage:               EquipmentUsageSectorAddressableMassStorage,
		SelectableEquipment: []string{"FSDISK", "RMDISK", "SCDISK"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "FFSYS",
		Usage:               EquipmentUsageSectorAddressableMassStorage,
		SelectableEquipment: []string{"FSDISK"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "FSCSI",
		Usage:               EquipmentUsageSectorAddressableMassStorage,
		SelectableEquipment: []string{"SCDISK"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "FRAM",
		Usage:               EquipmentUsageSectorAddressableMassStorage,
		SelectableEquipment: []string{"RMDISK"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "D",
		Usage:               EquipmentUsageWordAddressableMassStorage,
		SelectableEquipment: []string{"FSDISK", "RMDISK", "SCDISK"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "DFSYS",
		Usage:               EquipmentUsageWordAddressableMassStorage,
		SelectableEquipment: []string{"FSDISK"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "DSCSI",
		Usage:               EquipmentUsageWordAddressableMassStorage,
		SelectableEquipment: []string{"SCDISK"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "DRAM",
		Usage:               EquipmentUsageWordAddressableMassStorage,
		SelectableEquipment: []string{"RMDISK"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "T",
		Usage:               EquipmentUsageTape,
		SelectableEquipment: []string{"FSTAPE", "SCTAPE"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "TFSYS",
		Usage:               EquipmentUsageTape,
		SelectableEquipment: []string{"FSTAPE"},
	})
	cfg.importEquipmentEntry(&EquipmentEntry{
		Mnemonic:            "TSCSI",
		Usage:               EquipmentUsageTape,
		SelectableEquipment: []string{"SCTAPE"},
	})
	return cfg
}

func (cfg *Configuration) UpdateFromFile(fileName string) error {
	// TODO
	return fmt.Errorf("configuration UpdateFromFile() not yet implemented")
}

func (cfg *Configuration) importEquipmentEntry(entry *EquipmentEntry) {
	cfg.EquipmentTable[entry.Mnemonic] = entry
}
