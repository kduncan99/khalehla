// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
	"strings"
	"time"
)

/*
FS KEYIN - component DOES NOT EXIST, INPUT IGNORED
FS KEYIN - eqp-mnemonic EQUIPMENT MNEMONIC ILLEGAL, INPUT IGNORED
  (Exec) An equipment mnemonic cannot be entered on any keyin other than the DN,PACK keyin.
FS KEYIN - inhibits INHIBITS ILLEGAL, INPUT IGNORED
  (Exec) Inhibits cannot be entered on any keyin other than the MD keyin.
FS KEYIN NOT ALLOWED - DIRECTORY ID MUST BE dir-id
  (Exec) The directory-id of a pack and a directory-id specified on the keyin are opposite
  (for example, a local pack-id and shared were specified on the keyin).
FS KEYIN - NO UNIT EXISTS IN THE PREMOUNT ONLY STATUS, INPUT IGNORED
FS KEYIN - option OPTION DOES NOT EXIST, INPUT IGNORED
  (Exec) An illegal FS keyin was entered. The variable option is the option that you specified on your FS keyin.
FS NOT ALLOWED UNTIL MASS STORAGE INITIALIZED OR RECOVERED
  (Exec) An FS,PACK keyin is not allowed until the recovery files have been created or restored.
FS,PACK KEYIN ERROR - MHFS IS NOT AVAILABLE
  (Exec) The FS,PACK/SHARED keyin is not allowed because Multi-Host File Sharing (MHFS) is down or not available.
FS,PACK NOT ALLOWED - dir-id IS ILLEGAL DIRECTORY ID

Variations we accept:
	FS,[ CM | DISK | FDISK | MS | PACK | RDISK | TAPE ]
	FS node_name[,...]
	FS,ALL channel_name
*/

type FSKeyinHandler struct {
	exec            types.IExec
	source          types.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewFSKeyinHandler(exec types.IExec, source types.ConsoleIdentifier, options string, arguments string) types.KeyinHandler {
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
			return nodeMgr.IsValidNodeName(kh.arguments)
		}
		return len(kh.options) <= 6 && len(kh.arguments) == 0
	}

	split := strings.Split(kh.arguments, ",")
	if len(split) < 1 {
		return false
	}

	for _, name := range split {
		if !nodeMgr.IsValidNodeName(strings.ToUpper(name)) {
			return false
		}
	}
	return true
}

func (kh *FSKeyinHandler) GetCommand() string {
	return "FS"
}

func (kh *FSKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *FSKeyinHandler) GetArguments() string {
	return kh.arguments
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
	for sx := 0; sx < len(statStrings); {
		str := statStrings[sx]
		sx++
		if !strings.ContainsRune(str, '*') {
			if sx < len(statStrings) && !strings.ContainsRune(statStrings[sx], '*') {
				str = fmt.Sprintf("%-30s%s", str, statStrings[sx])
				sx++
			}
		}
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
	}
}

func (kh *FSKeyinHandler) getStatusStringForNode(nodeInfo types.NodeInfo) string {
	fm := kh.exec.GetFacilitiesManager()
	str := nodeInfo.GetNodeName() + " "
	str += nodeMgr.GetNodeStatusString(nodeInfo.GetNodeStatus(), nodeInfo.IsAccessible())
	if nodeInfo.GetNodeCategory() == nodeMgr.NodeCategoryDevice {
		devInfo := nodeInfo.(types.DeviceInfo)
		str += " " + fm.GetDeviceStatusDetail(devInfo.GetDeviceIdentifier())
	}
	return str
}

func (kh *FSKeyinHandler) handleAllForChannel() {
	nm := kh.exec.GetNodeManager()
	statStrings := make([]string, 0)
	chName := strings.ToUpper(kh.arguments)
	nodeInfo, err := nm.GetNodeInfoByName(chName)
	if err != nil {
		msg := fmt.Sprintf("FS KEYIN - %v DOES NOT EXIST, INPUT IGNORED", chName)
		kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
		return
	}
	statStrings = append(statStrings, kh.getStatusStringForNode(nodeInfo))

	if nodeInfo.GetNodeCategory() == nodeMgr.NodeCategoryChannel {
		chInfo := nodeInfo.(types.ChannelInfo)
		devInfos := chInfo.GetDeviceInfos()
		for _, di := range devInfos {
			statStrings = append(statStrings, kh.getStatusStringForNode(di))
		}
	}

	kh.emitStatusStrings(statStrings)
}

func (kh *FSKeyinHandler) handleAllOf(nodeCategory nodeMgr.NodeCategoryType, nodeType nodeMgr.NodeDeviceType) {
	nm := kh.exec.GetNodeManager()
	statStrings := make([]string, 0)
	if nodeCategory == nodeMgr.NodeCategoryChannel || nodeCategory == 0 {
		for _, chInfo := range nm.GetChannelInfos() {
			if nodeType == chInfo.GetNodeType() || nodeType == 0 {
				statStrings = append(statStrings, kh.getStatusStringForNode(chInfo))
			}
		}
	}

	if nodeCategory == nodeMgr.NodeCategoryDevice || nodeCategory == 0 {
		for _, devInfo := range nm.GetDeviceInfos() {
			if nodeType == devInfo.GetNodeType() || nodeType == 0 {
				statStrings = append(statStrings, kh.getStatusStringForNode(devInfo))
			}
		}
	}

	kh.emitStatusStrings(statStrings)
}

func (kh *FSKeyinHandler) handleComponentList() {
	nm := kh.exec.GetNodeManager()
	names := strings.Split(kh.arguments, ",")
	statStrings := make([]string, len(names))
	for nx, name := range names {
		ni, err := nm.GetNodeInfoByName(strings.ToUpper(name))
		if err != nil {
			msg := fmt.Sprintf("FS KEYIN - %v DOES NOT EXIST, INPUT IGNORED", name)
			kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
			return
		}

		statStrings[nx] = kh.getStatusStringForNode(ni)
	}

	kh.emitStatusStrings(statStrings)
}

func (kh *FSKeyinHandler) handleOption() {
	switch kh.options {
	case "ALL":
		kh.handleAllForChannel()
		return
	case "CM":
		kh.handleAllOf(nodeMgr.NodeCategoryChannel, 0)
		return
	case "DISK":
		kh.handleAllOf(nodeMgr.NodeCategoryDevice, nodeMgr.NodeDeviceDisk)
		return
	case "FDISK":
		// TODO
	case "MS":
		// TODO
	case "PACK":
		// TODO
	case "RDISK":
		// TODO
	case "TAPE":
		kh.handleAllOf(nodeMgr.NodeCategoryDevice, nodeMgr.NodeDeviceTape)
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
