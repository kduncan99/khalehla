// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

import (
	"io"
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

// DeviceReadyListener
// any entity which needs to be notified of devices going ready or not ready implements this
type DeviceReadyListener interface {
	NotifyDeviceReady(DeviceInfo, bool)
}

// type FacilitiesItem interface {
// 	GetInternalFileName() string
// 	GetFileName() string
// 	GetQualifier() string
// 	GetEquipmentCode() uint
// }

type IConsoleManager interface {
	CloseManager()
	Dump(destination io.Writer, indent string)
	InitializeManager() error // manager must stop the exec if it returns an error
	IsInitialized() bool
	ResetManager() error // manager must stop the exec if it returns an error
	SendReadOnlyMessage(message *ConsoleReadOnlyMessage)
	SendReadReplyMessage(message *ConsoleReadReplyMessage) error
	SendSystemMessages(message1 string, message2 string)
}

// IExec is the interface for the Exec, placed here to avoid package import cycles
type IExec interface {
	Close()
	Dump(destination io.Writer)
	GetConsoleManager() IConsoleManager
	GetFacilitiesManager() IFacilitiesManager
	GetKeyinManager() IKeyinManager
	GetNodeManager() INodeManager
	GetPhase() ExecPhase
	GetRunControlEntry() *RunControlEntry
	GetStopCode() StopCode
	GetStopFlag() bool
	HandleKeyIn(source ConsoleIdentifier, text string)
	SendExecReadOnlyMessage(message string)
	SendExecReadReplyMessage(message string, maxReplyChars int) (string, error)
	Stop(code StopCode)
}

type IFacilitiesManager interface {
	CloseManager()
	Dump(destination io.Writer, indent string)
	InitializeManager() error // manager must stop the exec if it returns an error
	IsInitialized() bool
	ResetManager() error // manager must stop the exec if it returns an error
	AssignDiskDeviceToExec(deviceId DeviceIdentifier) error
	GetDeviceStatusDetail(deviceId DeviceIdentifier) string
	NotifyDeviceReady(deviceInfo DeviceInfo, isReady bool)
}

type IKeyinManager interface {
	CloseManager()
	Dump(destination io.Writer, indent string)
	InitializeManager() error // manager must stop the exec if it returns an error
	IsInitialized() bool
	ResetManager() error // manager must stop the exec if it returns an error
	PostKeyin(source ConsoleIdentifier, text string)
}

type IMFDManager interface {
	CloseManager()
	Dump(destination io.Writer, indent string)
	InitializeManager() error // manager must stop the exec if it returns an error
	IsInitialized() bool
	ResetManager() error // manager must stop the exec if it returns an error
}

type INodeManager interface {
	CloseManager()
	Dump(destination io.Writer, indent string)
	InitializeManager() error // manager must stop the exec if it returns an error
	IsInitialized() bool
	ResetManager() error // manager must stop the exec if it returns an error
	GetChannelInfos() []ChannelInfo
	GetDeviceInfos() []DeviceInfo
	GetNodeInfoByName(nodeName string) (NodeInfo, error)
	GetNodeInfoByIdentifier(nodeId NodeIdentifier) (NodeInfo, error)
	RouteIo(ioPacket IoPacket)
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

// Manager is one of the top-level exec managers.
// They may have a goroutine operating for them.
type Manager interface {
	CloseManager()
	Dump(destination io.Writer, indent string)
	InitializeManager() error // manager must stop the exec if it returns an error
	IsInitialized() bool
	ResetManager() error // manager must stop the exec if it returns an error
}

// NodeInfo contains all the exec-managed information regarding a particular node
type NodeInfo interface {
	CreateNode()
	GetNodeCategory() NodeCategory
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	GetNodeStatus() NodeStatus
	GetNodeType() NodeType
	IsAccessible() bool
}
