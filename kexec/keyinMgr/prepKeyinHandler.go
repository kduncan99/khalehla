// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"khalehla/hardware"
	"khalehla/hardware/channels"
	"khalehla/hardware/ioPackets"
	"khalehla/kexec"
	"khalehla/kexec/facilitiesMgr"
	"khalehla/kexec/mfdMgr"
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
	"strconv"
	"strings"
	"time"
)

type PREPKeyinHandler struct {
	exec            kexec.IExec
	source          kexec.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time

	removable  bool
	deviceName string
	prepFactor int
	trackCount int
	packName   string
}

func NewPREPKeyinHandler(exec kexec.IExec, source kexec.ConsoleIdentifier, options string, arguments string) IKeyinHandler {
	return &PREPKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *PREPKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *PREPKeyinHandler) CheckSyntax() bool {
	// Syntax:
	//   PREP,[F|R] device,prepFactor,trackCount,pack_name
	split := strings.Split(kh.arguments, ",")
	if len(kh.options) != 1 || len(split) != 4 {
		return false
	}

	upOpts := strings.ToUpper(kh.options)
	if upOpts != "F" && upOpts != "R" {
		return false
	}
	kh.removable = upOpts == "R"

	var err error
	kh.deviceName = strings.ToUpper(split[0])
	kh.prepFactor, err = strconv.Atoi(split[1])
	if err != nil {
		return false
	}
	kh.trackCount, err = strconv.Atoi(split[2])
	if err != nil {
		return false
	}
	kh.packName = strings.ToUpper(split[3])
	if !kexec.IsValidNodeName(kh.deviceName) || !hardware.IsValidPackName(kh.packName) {
		return false
	}

	return true
}

func (kh *PREPKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *PREPKeyinHandler) GetCommand() string {
	return "PREP"
}

func (kh *PREPKeyinHandler) GetHelp() []string {
	return []string{
		"PREP,[F|R] device,prep_factor,track_count,pack_name",
		"Preps a virtual or real disk pack for use as mass storage",
		"Use F for fixed packs, R for removable packs",
		"prep_factor and track_count should be zero for real disk packs"}
}

func (kh *PREPKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *PREPKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *PREPKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *PREPKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *PREPKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *PREPKeyinHandler) process() {
	fm := kh.exec.GetFacilitiesManager().(*facilitiesMgr.FacilitiesManager)

	attr, ok := fm.GetNodeAttributesByName(kh.deviceName)
	if !ok {
		str := fmt.Sprintf("%v not found", kh.deviceName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	if attr.GetNodeCategoryType() != hardware.NodeCategoryDevice ||
		attr.GetNodeDeviceType() != hardware.NodeDeviceDisk {
		str := fmt.Sprintf("%v is not a disk device", kh.deviceName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	if !hardware.IsValidPrepFactor(hardware.PrepFactor(kh.prepFactor)) {
		kh.exec.SendExecReadOnlyMessage("Invalid PrepFactor", &kh.source)
		return
	}

	if kh.trackCount < 10000 || kh.trackCount > 262143 {
		kh.exec.SendExecReadOnlyMessage("Invalid TrackCount", &kh.source)
		return
	}

	// TODO Make sure the device is RV (after we get the RV keyin implemented)

	nodeId := attr.GetNodeIdentifier()
	nm := kh.exec.GetNodeManager().(*nodeMgr.NodeManager)
	nodeInfo, _ := nm.GetNodeInfoByIdentifier(nodeId)
	ddi := nodeInfo.(*nodeMgr.DiskDeviceInfo)
	dd := ddi.GetDiskDevice()

	blockSize, blockCount, currentPrepFactor, currentTrackCount := dd.GetDiskGeometry()
	if blockSize == 0 || blockCount == 0 || currentTrackCount == 0 {
		str := fmt.Sprintf("PREP %v pack is not properly formatted", kh.deviceName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
	} else {
		label := make([]pkg.Word36, currentPrepFactor)
		ioStat := kh.readBlock(nodeId, label, 0)
		if ioStat == ioPackets.IosInternalError {
			return
		} else if ioStat != ioPackets.IosComplete {
			// This is odd - the device knows the geometry, but we failed to read the disk label
			str := fmt.Sprintf("PREP %v IO error reading pack label", kh.deviceName)
			kh.exec.SendExecReadOnlyMessage(str, &kh.source)
			return
		}

		msg := fmt.Sprintf("PREP %v Relabel/Reprep Pack %v? Y/N", kh.deviceName, kh.packName)
		reply, err := kh.exec.SendExecRestrictedReadReplyMessage(msg, []string{"Y", "N"}, nil)
		if err != nil {
			return
		} else if reply == "N" {
			str := fmt.Sprintf("PREP %v canceled", kh.deviceName)
			kh.exec.SendExecReadOnlyMessage(str, &kh.source)
			return
		}
	}

	// Send an IofPrep to the disk to write the label and re-establish geometry
	pi := ioPackets.IoPrepInfo{
		PrepFactor:  hardware.PrepFactor(kh.prepFactor),
		TrackCount:  hardware.TrackCount(kh.trackCount),
		PackName:    kh.packName,
		IsRemovable: kh.removable,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeId,
		IoFunction:     ioPackets.IofPrep,
		PrepInfo:       &pi,
	}

	ioStat := kh.io(cp)
	if ioStat != ioPackets.IosComplete {
		str := fmt.Sprintf("PREP %v IO error %v writing pack label", kh.deviceName, ioStat)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	// Need to read the label for the next part
	label := make([]pkg.Word36, kh.prepFactor)
	ioStat = kh.readBlock(nodeId, label, 0)
	if ioStat == ioPackets.IosInternalError {
		str := fmt.Sprintf("PREP %v IO error %v (re)reading pack label", kh.deviceName, ioStat)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	// Populate and write the initial directory track
	dirTrack := make([]pkg.Word36, 1792)
	mfdMgr.PopulateInitialDirectoryTrack(label, !kh.removable, dirTrack)
	dirTrackDRWA := label[03].GetW()
	dirTrackAddr := dirTrackDRWA / 1792
	blocksPerTrack := label[04].GetH1()
	dirBlockAddr := dirTrackAddr * blocksPerTrack
	wordsPerBlock := label[04].GetH2()
	blockId := hardware.BlockId(dirBlockAddr)
	wx := uint64(0)
	for wx < 1792 {
		ioStat = kh.writeBlock(nodeId, dirTrack[wx:wx+wordsPerBlock], blockId)
		if ioStat == ioPackets.IosInternalError {
			str := fmt.Sprintf("PREP %v IO error %v writing directory track", kh.deviceName, ioStat)
			kh.exec.SendExecReadOnlyMessage(str, &kh.source)
			return
		}
		blockId++
		wx += wordsPerBlock
	}

	str := fmt.Sprintf("PREP %v Complete", kh.deviceName)
	kh.exec.SendExecReadOnlyMessage(str, &kh.source)
}

func (kh *PREPKeyinHandler) thread() {
	kh.threadStarted = true

	kh.process()

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}

// TODO when fac mgr is ready, we should assign the device to the exec and let fac mgr do the IOs.
func (kh *PREPKeyinHandler) io(
	cp *channels.ChannelProgram,
) ioPackets.IoStatus {
	kh.exec.GetNodeManager().RouteIo(cp)
	for cp.IoStatus == ioPackets.IosInProgress || cp.IoStatus == ioPackets.IosNotStarted {
		time.Sleep(10 * time.Millisecond)
	}

	return cp.IoStatus
}

func (kh *PREPKeyinHandler) readBlock(
	nodeId hardware.NodeIdentifier,
	buffer []pkg.Word36,
	blockId hardware.BlockId,
) ioPackets.IoStatus {
	cw := channels.ControlWord{
		Buffer:    buffer,
		Offset:    0,
		Length:    uint(len(buffer)),
		Direction: channels.DirectionForward,
		Format:    channels.TransferPacked,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeId,
		IoFunction:     ioPackets.IofRead,
		BlockId:        blockId,
		ControlWords:   []channels.ControlWord{cw},
	}

	return kh.io(cp)
}

func (kh *PREPKeyinHandler) writeBlock(
	nodeId hardware.NodeIdentifier,
	buffer []pkg.Word36,
	blockId hardware.BlockId,
) ioPackets.IoStatus {
	cw := channels.ControlWord{
		Buffer:    buffer,
		Offset:    0,
		Length:    uint(len(buffer)),
		Direction: channels.DirectionForward,
		Format:    channels.TransferPacked,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeId,
		IoFunction:     ioPackets.IofWrite,
		BlockId:        blockId,
		ControlWords:   []channels.ControlWord{cw},
	}

	return kh.io(cp)
}
