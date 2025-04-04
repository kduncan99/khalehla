// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"

	"khalehla/hardware"
	"khalehla/kexec"
	"khalehla/kexec/facilitiesMgr"
	hardware2 "khalehla/old/hardware"
	kexec2 "khalehla/old/kexec"
	nodeMgr2 "khalehla/old/kexec/nodeMgr"

	"strings"
	"time"
)

/*
Variations we accept:
	FS,[ CM | DISK | FDISK | MS | PACK | RDISK | TAPE ]
	FS node_name[,...]
	FS,ALL channel_name
*/

type FSKeyinHandler struct {
	exec    kexec.IExec
	source  kexec2.ConsoleIdentifier
	options string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewFSKeyinHandler(exec kexec.IExec, source kexec2.ConsoleIdentifier, options string, arguments string) IKeyinHandler {
	return &FSKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *FSKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *FSKeyinHandler) CheckSyntax() bool {
	if len(kh.options) != 0 {
		if kh.options == "ALL" {
			return kexec2.IsValidNodeName(kh.arguments)
		}
		return len(kh.options) <= 6 && len(kh.arguments) == 0
	}

	split := strings.Split(kh.arguments, ",")
	if len(split) < 1 {
		return false
	}

	for _, name := range split {
		if !kexec2.IsValidNodeName(strings.ToUpper(name)) {
			return false
		}
	}
	return true
}

func (kh *FSKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *FSKeyinHandler) GetCommand() string {
	return "FS"
}

func (kh *FSKeyinHandler) GetHelp() []string {
	return []string{
		"FS,[ CM | DISK[S] | FDISK | MS | PACK[S] | RDISK | TAPE[S] ]",
		"FS node_name[,...]",
		"FS,ALL channel_name",
		"Displays facility status for various system components"}
}

func (kh *FSKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *FSKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *FSKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *FSKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *FSKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *FSKeyinHandler) emitStatusStrings(statStrings []string) {
	if len(statStrings) == 0 {
		kh.exec.SendExecReadOnlyMessage("NO DEVICES", &kh.source)
	} else {
		for sx := 0; sx < len(statStrings); {
			str := statStrings[sx]
			sx++
			if !strings.ContainsRune(str, '*') {
				if sx < len(statStrings) && !strings.ContainsRune(statStrings[sx], '*') {
					str = fmt.Sprintf("%-33s%s", str, statStrings[sx])
					sx++
				}
			}
			kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		}
	}
}

func (kh *FSKeyinHandler) getStatusStringForNode(nodeId hardware.NodeIdentifier) string {
	fm := kh.exec.GetFacilitiesManager().(*facilitiesMgr.FacilitiesManager)
	attr, _ := fm.GetNodeAttributes(nodeId)
	return attr.GetNodeName() + " " + fm.GetNodeStatusString(nodeId)
}

func (kh *FSKeyinHandler) handleAllForChannel() {
	nm := kh.exec.GetNodeManager().(*nodeMgr2.NodeManager)
	statStrings := make([]string, 0)
	chName := strings.ToUpper(kh.arguments)
	nodeInfo, err := nm.GetNodeInfoByName(chName)
	if err != nil {
		msg := fmt.Sprintf("FS KEYIN - %v DOES NOT EXIST, INPUT IGNORED", chName)
		kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
		return
	}
	statStrings = append(statStrings, kh.getStatusStringForNode(nodeInfo.GetNodeIdentifier()))

	if nodeInfo.GetNodeCategoryType() == hardware.NodeCategoryChannel {
		chInfo := nodeInfo.(nodeMgr2.ChannelInfo)
		devInfos := chInfo.GetDeviceInfos()
		for _, di := range devInfos {
			statStrings = append(statStrings, kh.getStatusStringForNode(di.GetNodeIdentifier()))
		}
	}

	kh.emitStatusStrings(statStrings)
}

func (kh *FSKeyinHandler) handleAllOf(nodeCategory hardware.NodeCategoryType, nodeType hardware2.NodeDeviceType) {
	nm := kh.exec.GetNodeManager().(*nodeMgr2.NodeManager)
	statStrings := make([]string, 0)
	if nodeCategory == hardware.NodeCategoryChannel || nodeCategory == 0 {
		for _, chInfo := range nm.GetChannelInfos() {
			if nodeType == chInfo.GetNodeDeviceType() || nodeType == 0 {
				statStrings = append(statStrings, kh.getStatusStringForNode(chInfo.GetNodeIdentifier()))
			}
		}
	}

	if nodeCategory == hardware.NodeCategoryDevice || nodeCategory == 0 {
		for _, devInfo := range nm.GetDeviceInfos() {
			if nodeType == devInfo.GetNodeDeviceType() || nodeType == 0 {
				statStrings = append(statStrings, kh.getStatusStringForNode(devInfo.GetNodeIdentifier()))
			}
		}
	}

	kh.emitStatusStrings(statStrings)
}

func (kh *FSKeyinHandler) handleComponentList() {
	fm := kh.exec.GetFacilitiesManager().(*facilitiesMgr.FacilitiesManager)
	names := strings.Split(kh.arguments, ",")
	statStrings := make([]string, len(names))
	for nx, name := range names {
		attr, ok := fm.GetNodeAttributesByName(name)
		if !ok {
			msg := fmt.Sprintf("FS KEYIN - %v DOES NOT EXIST, INPUT IGNORED", name)
			kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
			return
		}

		statStrings[nx] = kh.getStatusStringForNode(attr.GetNodeIdentifier())
	}

	kh.emitStatusStrings(statStrings)
}

func (kh *FSKeyinHandler) handleOption() {
	switch kh.options {
	case "ALL":
		kh.handleAllForChannel()
		return
	case "CM":
		kh.handleAllOf(hardware.NodeCategoryChannel, 0)
		return
	case "DISK":
		kh.handleAllOf(hardware.NodeCategoryDevice, hardware2.NodeDeviceDisk)
		return
	case "DISKS":
		kh.handleAllOf(hardware.NodeCategoryDevice, hardware2.NodeDeviceDisk)
		return
	case "FDISK":
		// TODO FS,FDISK
		// "NO FIXED DISK CONFIGURED"
	case "MS":
		// TODO FS,MS
	case "PACK":
		// TODO FS,PACK
	case "PACKS":
		// TODO FS,PACKS
	case "RDISK":
		// TODO FS,RDISK
		// "NO REMOVABLE DISKS PRESENT"
	case "TAPE":
		kh.handleAllOf(hardware.NodeCategoryDevice, hardware2.NodeDeviceTape)
		return
	case "TAPES":
		kh.handleAllOf(hardware.NodeCategoryDevice, hardware2.NodeDeviceTape)
		return
	}

	msg := fmt.Sprintf("FS KEYIN - %v OPTION DOES NOT EXIST, INPUT IGNORED", kh.options)
	kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
}

func (kh *FSKeyinHandler) thread() {
	kh.threadStarted = true

	if len(kh.options) > 0 {
		kh.handleOption()
	} else {
		kh.handleComponentList()
	}

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
