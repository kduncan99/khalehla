// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

import (
	"io"
	"khalehla/kexec/config"
	"time"
)

// Channel manages async communication with the various deviceInfos assigned to it.
// It may also manage caching, automatic mounting, or any other various activities
// on behalf of the exec.
type Channel interface {
	AssignDevice(deviceIdentifier DeviceIdentifier, device Device) error
	Dump(destination io.Writer, indent string)
	GetNodeType() NodeType
	StartIo(ioPacket IoPacket)
}

// ChannelInfo is intended primarily as a means of documenting the use of a more generic NodeInfo
type ChannelInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetChannel() Channel
	GetChannelName() string
	GetChannelIdentifier() ChannelIdentifier
	GetDeviceInfos() []DeviceInfo
	GetNodeCategory() NodeCategory
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	GetNodeStatus() NodeStatus
	GetNodeType() NodeType
	IsAccessible() bool
}

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

// Device manages real or pseudo IO operations for a particular virtual device.
// It may do so synchronously or asynchronously
type Device interface {
	Dump(destination io.Writer, indent string)
	GetNodeType() NodeType
	IsReady() bool
	StartIo(ioPacket IoPacket)
}

// DeviceInfo is intended primarily as a means of documenting the use of a more generic NodeInfo
type DeviceInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetChannelInfos() []ChannelInfo
	GetDevice() Device
	GetDeviceIdentifier() DeviceIdentifier
	GetDeviceName() string
	GetNodeCategory() NodeCategory
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	GetNodeStatus() NodeStatus
	GetNodeType() NodeType
	IsAccessible() bool
	IsReady() bool
	SetIsAccessible(bool)
	SetIsReady(flag bool)
}

type FacilitiesItem interface {
	GetInternalFileName() string
	GetFileName() string
	GetQualifier() string
	GetEquipmentCode() uint
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
	NotifyDeviceReady(deviceInfo DeviceInfo, isReady bool)
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
	GetChannelInfos() []ChannelInfo
	GetDeviceInfos() []DeviceInfo
	GetNodeInfoByName(nodeName string) (NodeInfo, error)
	GetNodeInfoByIdentifier(nodeId NodeIdentifier) (NodeInfo, error)
	RouteIo(ioPacket IoPacket)
	SetNodeStatus(nodeId NodeIdentifier, status NodeStatus) error
}

// IoPacket contains all the information necessary for a Channel to route an IO operation,
// and for a device to perform that IO operation.
type IoPacket interface {
	GetDeviceIdentifier() DeviceIdentifier
	GetNodeType() NodeType
	GetIoFunction() IoFunction
	GetIoStatus() IoStatus
	SetIoStatus(ioStatus IoStatus)
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

// NodeInfo contains all the exec-managed information regarding a particular node
type NodeInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetNodeCategory() NodeCategory
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	GetNodeStatus() NodeStatus
	GetNodeType() NodeType
	IsAccessible() bool
}
