// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

import (
	"io"
	"khalehla/kexec/config"
	"khalehla/kexec/nodeMgr"
	"time"
)

// Console is a unit which actually acts as an operating system console endpoint.
// One example is the StandardConsole which is always around.
// One might also implement a DemandConsole for RSI @@CONS, or a net-based console for a web browser.
type Console interface {
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
	HandleKeyIn(source ConsoleIdentifier, text string)
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
	AssignDiskDeviceToExec(deviceId DeviceIdentifier) error
	GetDeviceStatusDetail(deviceId DeviceIdentifier) string
	GetDiskAttributes(deviceId DeviceIdentifier) (*DiskAttributes, error)
	IsDeviceAssigned(deviceId DeviceIdentifier) bool
	NotifyDeviceReady(deviceInfo nodeMgr.DeviceInfo, isReady bool)
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
	//	NotifyDeviceReady(deviceInfo DeviceInfo, isReady bool)
}

type INodeManager interface {
	Boot() error // invoked when the exec is booting - return an error to stop the boot
	Close()      // invoked when the application is shutting down
	Dump(destination io.Writer, indent string)
	Initialize() error // invoked when the application is starting up
	Stop()             // invoked when the exec is stopping
	GetChannelInfos() []nodeMgr.ChannelInfo
	GetDeviceInfos() []nodeMgr.DeviceInfo
	GetNodeInfoByName(nodeName string) (nodeMgr.NodeInfo, error)
	GetNodeInfoByIdentifier(nodeId NodeIdentifier) (nodeMgr.NodeInfo, error)
	RouteIo(ioPacket nodeMgr.IoPacket)
	SetNodeStatus(nodeId NodeIdentifier, status NodeStatus) error
}

type KeyinHandler interface {
	Abort()
	CheckSyntax() bool
	GetCommand() string
	GetOptions() string
	GetArguments() string
	GetTimeFinished() time.Time
	Invoke()
	IsAllowed() bool
	IsDone() bool
}
