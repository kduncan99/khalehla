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
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
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
}

func NewPREPKeyinHandler(exec kexec.IExec, source kexec.ConsoleIdentifier, options string, arguments string) KeyinHandler {
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
	//   PREP,[F|R] device,pack_name

	if len(kh.options) != 1 || len(kh.arguments) != 2 {
		return false
	}

	if kh.options != "F" && kh.options != "R" {
		return false
	}

	return true
}

func (kh *PREPKeyinHandler) GetCommand() string {
	return "PREP"
}

func (kh *PREPKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *PREPKeyinHandler) GetArguments() string {
	return kh.arguments
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

	remFlag := kh.options == "R"

	split := strings.Split(kh.arguments, ",")
	deviceName := strings.ToUpper(split[0])
	packName := strings.ToUpper(split[1])

	attr, ok := fm.GetNodeAttributesByName(deviceName)
	if !ok {
		str := fmt.Sprintf("%v not found", deviceName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	if attr.GetNodeCategoryType() != hardware.NodeCategoryDevice ||
		attr.GetNodeDeviceType() != hardware.NodeDeviceDisk {
		str := fmt.Sprintf("%v is not a disk device", deviceName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	// TODO Make sure the device is RV (after we get the RV keyin implemented)

	if !hardware.IsValidPackName(packName) {
		str := fmt.Sprintf("PREP %v %v is not a valid pack name", deviceName, packName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	nodeId := attr.GetNodeIdentifier()
	nm := kh.exec.GetNodeManager().(*nodeMgr.NodeManager)
	nodeInfo, _ := nm.GetNodeInfoByIdentifier(nodeId)
	ddi := nodeInfo.(*nodeMgr.DiskDeviceInfo)
	dd := ddi.GetDiskDevice()

	blockSize, blockCount, trackCount := dd.GetDiskGeometry()
	if blockSize == 0 || blockCount == 0 || trackCount == 0 {
		str := fmt.Sprintf("PREP %v pack is not properly formatted", deviceName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	prepFactor := hardware.PrepFactorFromBlockSize[blockSize]
	label := make([]pkg.Word36, prepFactor)
	ioStat := kh.readBlock(nodeId, label, 0)
	if ioStat == ioPackets.IosInternalError {
		return
	} else if ioStat != ioPackets.IosComplete {
		// This is odd - the device knows the geometry, but we failed to read the disk label
		str := fmt.Sprintf("PREP %v IO error reading pack label", deviceName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	// TODO console R/Rply message display current geometry, ask if we want to overwrite it

	// Send an IofPrep to the disk to write the label and re-establish geometry
	cp := &channels.ChannelProgram{
		NodeIdentifier:   nodeId,
		IoFunction:       ioPackets.IofPrep,
		// PrepInfo
	}
	kh.exec.GetNodeManager().RouteIo(cp)
	for cp.IoStatus == ioPackets.IosComplete || cp.IoStatus == ioPackets.IosNotStarted {
		time.Sleep(10 * time.Millisecond)
	}
	if cp.IoStatus != ioPackets.IosComplete {
		str := fmt.Sprintf("PREP %v IO error writing pack label", deviceName)
		kh.exec.SendExecReadOnlyMessage(str, &kh.source)
		return
	}

	// write initial directory track...
	// actually, we only need to write sectors 0 and 1, and any slop necessary to pad out to the prep factor.
	blocksPerTrack := hardware.BlockCount(1792 / prepFactor)
	dirTrack := make([]pkg.Word36, 1792)
	availableTracks := trackCount - 2 // subtract label track and first directory track

	// sector 0
	das := dirTrack[0:28]
	das[1].SetW(0_600000_000000) // first 2 sectors are allocated
	for dx := 3; dx < 27; dx += 3 {
		das[dx].SetW(0_400000_000000)
	}
	das[27].SetW(0_400000_000000)

	// sector 1
	s1 := dirTrack[28:56]
	// leave +0 and +1 alone (We aren't doing HMBT/SMBT so we don't need the addresses)
	s1[2].SetW(uint64(availableTracks))
	s1[3].SetW(uint64(availableTracks))
	s1[4].FromStringToFieldata(packName)
	if !remFlag {
		s1[5].SetH1(0_400000)
	}
	s1[010].SetT1(uint64(blocksPerTrack))
	s1[010].SetS3(1) // Sector 1 version
	s1[010].SetT3(uint64(prepFactor))

	dirBlockId := hardware.BlockId(blocksPerTrack) // assuming directory track is the second logical track
	for wx := 0; wx < 56; wx += int(prepFactor) {
		subBuffer := dirTrack[wx : wx+int(prepFactor)]
		ioStat = kh.writeBlock(nodeId, subBuffer, dirBlockId)
		if ioStat == ioPackets.IosInternalError {
			return
		}
		dirBlockId++
	}
}

func (kh *PREPKeyinHandler) thread() {
	kh.threadStarted = true

	kh.process()

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}

// TODO when fac mgr is ready, we should assign the device to the exec and let fac mgr do the IOs.
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

	kh.exec.GetNodeManager().RouteIo(cp)
	for cp.IoStatus == ioPackets.IosComplete || cp.IoStatus == ioPackets.IosNotStarted {
		time.Sleep(10 * time.Millisecond)
	}

	return cp.IoStatus
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

	kh.exec.GetNodeManager().RouteIo(cp)
	for cp.IoStatus == ioPackets.IosComplete || cp.IoStatus == ioPackets.IosNotStarted {
		time.Sleep(10 * time.Millisecond)
	}

	return cp.IoStatus
}
