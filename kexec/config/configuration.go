// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package config

var defaultMasterAccountId = ""
var defaultOverheadAccountId = "INSTALLATION"
var defaultOverheadUserId = "INSTALLATION"
var defaultPrivilegedAccountId = "123456"
var defaultSecurityOfficerUserId = ""
var defaultSystemAccountId = "SYSTEM"
var defaultSystemProjectId = "SYSTEM"
var defaultSystemQualifier = "SYS$"
var defaultSystemRunId = "EXEC-8"
var defaultSystemUserId = "EXEC-8"

type Configuration struct {
	MasterAccountId     string
	OverheadAccountId   string
	OverheadUserId      string
	PrivilegedAccountId string
	SecurityOfficerId   string
	SystemAccountId     string
	SystemProjectId     string
	SystemQualifier     string
	SystemRunId         string
	SystemUserId        string
}

func NewConfiguration() *Configuration {
	cfg := &Configuration{}

	cfg.MasterAccountId = defaultMasterAccountId
	cfg.OverheadAccountId = defaultOverheadAccountId
	cfg.OverheadUserId = defaultOverheadUserId
	cfg.PrivilegedAccountId = defaultPrivilegedAccountId
	cfg.SecurityOfficerId = defaultSecurityOfficerUserId
	cfg.SystemAccountId = defaultSystemAccountId
	cfg.SystemProjectId = defaultSystemProjectId
	cfg.SystemQualifier = defaultSystemQualifier
	cfg.SystemRunId = defaultSystemRunId
	cfg.SystemUserId = defaultSystemUserId

	return cfg
}
