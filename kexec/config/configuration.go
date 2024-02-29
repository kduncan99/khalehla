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
}

func NewConfiguration() *Configuration {
	cfg := &Configuration{}

	cfg.AssignMnemonic = defaultAssignMnemonic
	cfg.LogIOs = defaultLogIOs
	cfg.MasterAccountId = defaultMasterAccountId
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
