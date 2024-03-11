// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"io"
	"khalehla/hardware"
	"khalehla/kexec/config"
)

// IConsole is a unit which actually acts as an operating system console endpoint.
// One example is the StandardConsole which is always around.
// One might also implement a DemandConsole for RSI @@CONS, or a net-based console for a web browser.
type IConsole interface {
	ClearReadReplyMessage(messageId int) error
	Close()
	Dump(destination io.Writer, indent string)
	PollSolicitedInput() (*string, int, error)
	PollUnsolicitedInput() (*string, error)
	IsConnected() bool
	Reset() error
	SendReadOnlyMessage(text string) error
	SendSystemMessages(text1 string, text2 string) error
	SendReadReplyMessage(text string, maxChars int) (int, error)
}

type IConsoleManager interface {
	Boot() error // invoked when the exec is booting - return an error to stop the boot
	Close()      // invoked when the application is shutting down
	Dump(destination io.Writer, indent string)
	Initialize() error // invoked when the application is starting up
	Stop()             // invoked when the exec is stopping
	SendReadOnlyMessage(message *ConsoleReadOnlyMessage)
	SendReadReplyMessage(message *ConsoleReadReplyMessage) error
	SendSystemMessages(message1 string, message2 string)
}

// IExec is the interface for the Exec, placed here to avoid package import cycles
type IExec interface {
	Boot(session uint, jumpKeys []bool, invokerChannel chan StopCode)
	Close()
	GetConfiguration() *config.Configuration
	GetConsoleManager() IConsoleManager
	GetFacilitiesManager() IFacilitiesManager
	GetJumpKey(jkNumber int) bool
	GetKeyinManager() IKeyinManager
	GetMFDManager() IMFDManager
	GetNodeManager() INodeManager
	GetPhase() ExecPhase
	GetRunControlEntry() *RunControlEntry
	GetStopCode() StopCode
	GetStopFlag() bool
	Initialize() error
	PerformDump(fullFlag bool) (string, error)
	SendExecReadOnlyMessage(message string, routing *ConsoleIdentifier)
	SendExecReadReplyMessage(message string, maxReplyChars int) (string, error)
	SendExecRestrictedReadReplyMessage(message string, accepted []string) (string, error)
	SetJumpKey(jkNumber int, value bool)
	Stop(code StopCode)
}

type IFacilitiesManager interface {
	Boot() error // invoked when the exec is booting - return an error to stop the boot
	Close()      // invoked when the application is shutting down
	Dump(destination io.Writer, indent string)
	Initialize() error // invoked when the application is starting up
	Stop()             // invoked when the exec is stopping
	//AssignDiskDeviceToExec(nodeId NodeIdentifier) error
	GetDiskAttributes(identifier hardware.NodeIdentifier) (*DiskAttributes, bool)
	GetNodeAttributes(nodeId hardware.NodeIdentifier) (NodeAttributes, bool)
	//IsDeviceAssigned(nodeId NodeIdentifier) bool
	//NotifyDeviceReady(nodeId NodeIdentifier, isReady bool)
	SetNodeStatus(nodeId hardware.NodeIdentifier, status FacNodeStatus) error
}

type IKeyinManager interface {
	Boot() error // invoked when the exec is booting - return an error to stop the boot
	Close()      // invoked when the application is shutting down
	Dump(destination io.Writer, indent string)
	Initialize() error // invoked when the application is starting up
	Stop()             // invoked when the exec is stopping
	PostKeyin(source ConsoleIdentifier, text string)
}

// IManager is one of the top-level exec managers.
// They may have a goroutine operating for them.
type IManager interface {
	Boot() error // invoked when the exec is booting - return an error to stop the boot
	Close()      // invoked when the application is shutting down
	Dump(destination io.Writer, indent string)
	Initialize() error // invoked when the application is starting up
	Stop()             // invoked when the exec is stopping
}

type IMFDManager interface {
	Boot() error // invoked when the exec is booting - return an error to stop the boot
	Close()      // invoked when the application is shutting down
	Dump(destination io.Writer, indent string)
	Initialize() error // invoked when the application is starting up
	InitializeMassStorage()
	RecoverMassStorage()
	Stop() // invoked when the exec is stopping
}

type INodeManager interface {
	Boot() error // invoked when the exec is booting - return an error to stop the boot
	Close()      // invoked when the application is shutting down
	Dump(destination io.Writer, indent string)
	Initialize() error // invoked when the application is starting up
	Stop()             // invoked when the exec is stopping
}
