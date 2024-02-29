// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package config

import (
	"fmt"
)

var defaultAssignMnemonic = "F"
var defaultLogIOs = false
var defaultMasterAccountId = ""
var defaultMaxGranules = uint64(256)
var defaultOverheadAccountId = "INSTALLATION"
var defaultOverheadUserId = "INSTALLATION"
var defaultPrivilegedAccountId = "123456"
var defaultSecurityOfficerUserId = ""
var defaultSystemAccountId = "SYSTEM"
var defaultSystemProjectId = "SYSTEM"
var defaultSystemQualifier = "SYS$"
var defaultSystemReadKey = "RDKEY"
var defaultSystemRunId = "EXEC-8"
var defaultSystemUserId = "EXEC-8"
var defaultSystemWriteKey = "WRKEY"

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
	AssignMnemonic      string
	LogIOs              bool
	MasterAccountId     string
	MaxGranules         uint64
	OverheadAccountId   string
	OverheadUserId      string
	PrivilegedAccountId string
	SecurityOfficerId   string
	SystemAccountId     string
	SystemProjectId     string
	SystemQualifier     string
	SystemReadKey       string
	SystemRunId         string
	SystemUserId        string
	SystemWriteKey      string
	EquipmentTable      map[string]*EquipmentEntry // key is mnemonic

	// TODO -- and this is not exhaustive...
	// RESDUCLR (residue_clear) zeroes tracks when allocating them
	// ACCTINTRES SYS$*ACCOUNT$R1 SYS$*SEC@ACCTINFO files initial reserve
	// ACCTASGMNE (ditto) assign mnemonic
	// DLOCASGMNE SYS$*DLOC$ assign mnemonic
	// GENFASGMNE SYS$*GENF$ assign mnemonic
	// GENFINTRES (ditto) initial reserve
	// LIBASGMNE SYS$*LIB$ assign mnemonic
	// LIBINTRES (ditto) initial reserve [0]
	// LIBMAXSIZ (ditto) max size [99999]
	// MAXCRD max cards [100]
	// MAXGRN already done
	// MAXPAG max pages [100]
	// MDFALT deafult sector-addressable mnemonic ['F']
	// MSTRACC master account [''] when blank and ACCOUNT$R1 is being initialized, prompt operator
	// RELUNUSEDREM release unused allocated space on removable packs at @FREE
	// RELUNUSEDRES release unused initial reserve at @FREE for fixed
	// RUNASGMNE SYS$*RUN$ assign mnemonic
	// RUNINTRES (ditto) initial reserve [1]
	// RUNMAXSIZE (ditto) [20000]
	// SACRDASGMNE SYS$*SEC@ACR$ assign mnemonic
	// SACRDINTRES (ditto) initial reserve [2]
	// SECOFFDEF security officer userid
	// SMDTFASGMNE SYS$*SMDTF$ (?) assign mnemonic
	// SMDTFINTRES (ditto) initial reserve [1]
	// SSEQPT default mnemonic for exec tape requests ['T']
	// SSPBP files are private by account [true]
	// SYMFBUF words used for standard/alt read/write buffers [224]
	// TDFALT default mnemonic for user tape requests ['T']
	// TLAUTO automatic tape labeling [false]
	// TPFMAXSIZ initial max size of TPF$
	// TPFTYP equip type for TPF$ ['F']
	// TPOWN tape access is restricted by account [false]
	// TPQUAL qualifier for TIP/Exec files ['TIP$']
	// TPRKEY TIP/Exec file read key ['++++++']
	// TPWKEY TIP/Exec file write key ['++++++']
	// TRMXCO terminate runs on max cards [false]
	// TRMXPO terminate runs on max pages [false]
	// TRMXT terminate runs on max time [false]
	// USERASGMNE SYS$*SEC@USERID$ assign mnemonic
	// USERINTRES (ditto) initial reserve [1]
	// WDFALT default mnemonic for word-addressable requests ['D']
}

func NewConfiguration() *Configuration {
	cfg := &Configuration{}

	cfg.AssignMnemonic = defaultAssignMnemonic
	cfg.LogIOs = defaultLogIOs
	cfg.MasterAccountId = defaultMasterAccountId
	cfg.MaxGranules = defaultMaxGranules
	cfg.OverheadAccountId = defaultOverheadAccountId
	cfg.OverheadUserId = defaultOverheadUserId
	cfg.PrivilegedAccountId = defaultPrivilegedAccountId
	cfg.SecurityOfficerId = defaultSecurityOfficerUserId
	cfg.SystemAccountId = defaultSystemAccountId
	cfg.SystemProjectId = defaultSystemProjectId
	cfg.SystemQualifier = defaultSystemQualifier
	cfg.SystemReadKey = defaultSystemReadKey
	cfg.SystemRunId = defaultSystemRunId
	cfg.SystemUserId = defaultSystemUserId
	cfg.SystemWriteKey = defaultSystemWriteKey
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
