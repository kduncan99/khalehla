// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package config

import (
	"fmt"
)

var defaultMasterAccount = ""
var defaultMaxGranules = 256
var defaultSecurityOfficerUserId = ""

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
	//AssignMnemonic      string
	//LogIOs              bool
	//MasterAccountId     string
	//MaxGranules         uint64
	//OverheadAccountId   string
	//OverheadUserId      string
	//PrivilegedAccountId string
	//SecurityOfficerId   string
	//SystemAccountId     string
	//SystemProjectId     string
	//SystemQualifier     string
	//SystemReadKey       string
	//SystemRunId         string
	//SystemUserId        string
	//SystemWriteKey      string
	MasterAccountId       string                     // could be empty, in which case operator is prompted when ACCOUNT$R1 is created
	MaxGranules           uint64                     // max granules if not specified on @ASG or @CAT
	SecurityOfficerUserId string                     // could be empty, in which case operator is prompted at boot time
	EquipmentTable        map[string]*EquipmentEntry // key is mnemonic

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
	// MAXGRN max granules if not specified
	// MAXPAG max pages [100]
	// MDFALT deafult sector-addressable mnemonic ['F']
	// MSTRACC master account [''] when blank and ACCOUNT$R1 is being initialized, prompt operator
	// OVRACC account ID for overhead runs (e.g., SYS, ROLBAK, ROLOUT)
	// OVRUSR user ID for overhead runs
	// PRIVAC '123456' privileged account number prevents tape label blocks from being read
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

	cfg.MasterAccountId = defaultMasterAccount
	cfg.SecurityOfficerUserId = defaultSecurityOfficerUserId

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
