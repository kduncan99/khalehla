// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package packUtil

import (
	"errors"
	"fmt"
	"khalehla/hardware"
	"khalehla/hardware/channels"
	"khalehla/hardware/devices"
	"khalehla/hardware/ioPackets"
	"khalehla/kexec/mfdMgr"
	"khalehla/pkg"
	"os"
	"strconv"
	"time"
)

func DoUsage() {
	fmt.Println("Usage:")
	fmt.Println("    packUtil prep {file_name} {pack_name} {prep_factor} {track_count} [ REM ]")
	fmt.Println("    packUtil show {file_name}")
}

func DoPrep(args []string) error {
	if len(args) < 4 || len(args) > 5 {
		return fmt.Errorf("incorrect number of arguments for prep command")
	}

	fileName := args[0]

	packName := args[1]
	if !hardware.IsValidPackName(packName) {
		return fmt.Errorf("invalid pack name")
	}

	pfInt, err := strconv.Atoi(args[2])
	if err != nil || pfInt <= 0 {
		return fmt.Errorf("error in prepfactor argument")
	}
	prepFactor := hardware.PrepFactor(pfInt)
	if !hardware.IsValidPrepFactor(prepFactor) {
		return fmt.Errorf("invalid prep factor (use 28, 56, 112, 224, 448, 896, or 1792)")
	}

	trackCount, err := strconv.Atoi(args[3])
	if err != nil || trackCount <= 0 {
		return fmt.Errorf("error in trackCount argument")
	}
	if trackCount < 10000 {
		return fmt.Errorf("invalid track count - must be at least 10000")
	}

	removable := false
	if len(args) == 5 {
		if args[4] != "REM" {
			return fmt.Errorf("optional argument for prep command is not REM")
		}
		removable = true
	}

	dc := channels.NewDiskChannel()
	dd := devices.NewFileSystemDiskDevice(nil)
	nodeId := hardware.NodeIdentifier(pkg.NewFromStringToFieldata("DISK0", 1)[0])
	_ = dc.AssignDevice(nodeId, dd)

	err = ioMount(dc, nodeId, fileName)
	if err != nil {
		return err
	}

	err = ioPrep(dc, nodeId, prepFactor, hardware.TrackCount(trackCount), packName, removable)
	if err != nil {
		return err
	}

	// Need to read the label for the next part
	label := make([]pkg.Word36, prepFactor)
	err = ioRead(dc, nodeId, label, 0, prepFactor)
	if err != nil {
		return err
	}

	// Populate and write the initial directory track
	dirTrack := make([]pkg.Word36, 1792)
	mfdMgr.PopulateInitialDirectoryTrack(label, !removable, dirTrack)
	dirTrackDRWA := label[03].GetW()
	dirTrackAddr := dirTrackDRWA / 1792
	blocksPerTrack := label[04].GetH1()
	dirBlockAddr := dirTrackAddr * blocksPerTrack
	blockId := hardware.BlockId(dirBlockAddr)
	wx := uint(0)
	for wx < 1792 {
		err = ioWrite(dc, nodeId, dirTrack[wx:wx+uint(prepFactor)], blockId, prepFactor)
		if err != nil {
			return err
		}
		blockId++
		wx += uint(prepFactor)
	}

	err = showLabelAndDirectory(dc, nodeId, prepFactor, true)
	if err != nil {
		return err
	}

	err = ioUnmount(dc, nodeId)
	if err != nil {
		return err
	}

	return nil
}

func DoShow(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments for show command")
	}

	fileName := args[0]
	if _, err := os.Stat(fileName); err == nil {
		// skip on down
	} else if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file %v does not exist", fileName)
	} else {
		return fmt.Errorf("cannot open file %v:%v", fileName, err)
	}

	dc := channels.NewDiskChannel()
	dd := devices.NewFileSystemDiskDevice(nil)
	nodeId := hardware.NodeIdentifier(1)
	_ = dc.AssignDevice(nodeId, dd)

	err := ioMount(dc, nodeId, fileName)
	if err != nil {
		return err
	}

	_, _, prepFactor, _ := dd.GetDiskGeometry()
	if prepFactor == 0 {
		return fmt.Errorf("pack is not prepped")
	}

	err = showLabelAndDirectory(dc, nodeId, prepFactor, true)
	if err != nil {
		return err
	}

	err = ioUnmount(dc, nodeId)
	if err != nil {
		return err
	}

	return nil
}

func showLabelAndDirectory(
	channel *channels.DiskChannel,
	devId hardware.NodeIdentifier,
	prepFactor hardware.PrepFactor,
	interpret bool,
) error {
	label := make([]pkg.Word36, prepFactor)
	err := ioRead(channel, devId, label, 0, prepFactor)
	if err != nil {
		return err
	}

	fmt.Println("Label Record:")
	pkg.DumpWord36Buffer(label[0:28], 7)
	if interpret {
		fmt.Printf("Pack Name:            %v%v\n", label[1].ToStringAsAscii(), label[2].ToStringAsAscii())
		fmt.Printf("First Dir Track DRWA: %v\n", label[3].ToStringAsOctal())
		fmt.Printf("Records Per Track:    %d\n", label[4].GetH1())
		fmt.Printf("Words Per Record:     %d\n", label[4].GetH2())
		fmt.Printf("VOL1 Version:         %d\n", label[014].GetS2())
		fmt.Printf("Disk Capacity:        %d tracks\n", label[016].GetW())
		fmt.Printf("Words Per Phys Record:%d\n", label[017].GetH1())
		fmt.Printf("Total Blocks:         %d\n", label[021].GetW())
	}

	dirTrackAddr := label[3].GetW()
	blocksPerTrack := 1792 / prepFactor
	dirBlockId := hardware.BlockId(uint(dirTrackAddr) / 1792 * uint(blocksPerTrack))
	dirBlock := make([]pkg.Word36, prepFactor)
	err = ioRead(channel, devId, dirBlock, dirBlockId, prepFactor)
	if err != nil {
		return err
	}

	fmt.Println("Sector 0:")
	pkg.DumpWord36Buffer(dirBlock[0:28], 7)

	var sector1 []pkg.Word36
	if prepFactor == 28 {
		dirBlockId++
		sector1 = make([]pkg.Word36, prepFactor)
		err = ioRead(channel, devId, sector1, dirBlockId, prepFactor)
		if err != nil {
			return err
		}
	} else {
		sector1 = dirBlock[28:56]
	}

	fmt.Println("Sector 1:")
	pkg.DumpWord36Buffer(sector1, 7)

	if interpret {
		fmt.Printf("HMBT DRWAddr:     %v\n", sector1[0].ToStringAsOctal())
		fmt.Printf("SMBT DRWAddr:     %v\n", sector1[1].ToStringAsOctal())
		fmt.Printf("Max Avail Tracks: %012o (%v)\n", sector1[2].GetW(), sector1[2].GetW())
		fmt.Printf("Available Tracks: %012o (%v)\n", sector1[3].GetW(), sector1[3].GetW())
		fmt.Printf("PackId:           %v\n", sector1[4].ToStringAsFieldata())
		remStr := ""
		if sector1[5].GetH1() == 0 {
			remStr = " (removable)"
		}
		fmt.Printf("LDAT Index:       %06o%v\n", sector1[5].GetH1(), remStr)
		fmt.Printf("MBT Length:       %06o\n", sector1[5].GetH2())
		fmt.Printf("Records/Track:    %v\n", sector1[010].GetT1())
		fmt.Printf("Words/Record:     %v\n", sector1[010].GetT3())
		fmt.Printf("Removable LDAT:   %06o\n", sector1[020].GetH1())
		fmt.Printf("DAS Offset:       %06o\n", sector1[020].GetH2())
	}

	return nil
}

func io(
	ch *channels.DiskChannel,
	cp *channels.ChannelProgram,
) error {
	ch.StartIo(cp)
	for cp.IoStatus == ioPackets.IosNotStarted || cp.IoStatus == ioPackets.IosInProgress {
		time.Sleep(10 * time.Millisecond)
	}
	if cp.IoStatus != ioPackets.IosComplete {
		fmt.Printf("IO error:%v status:%v\n", cp.GetString(), ioPackets.IoStatusTable[cp.IoStatus])
		return fmt.Errorf("error:%v", cp.GetString())
	}

	return nil
}

func ioMount(
	ch *channels.DiskChannel,
	nodeIdentifier hardware.NodeIdentifier,
	fileName string,
) error {
	mi := &ioPackets.IoMountInfo{
		Filename:     fileName,
		WriteProtect: false,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeIdentifier,
		IoFunction:     ioPackets.IofMount,
		MountInfo:      mi,
	}
	return io(ch, cp)
}

func ioPrep(
	ch *channels.DiskChannel,
	nodeIdentifier hardware.NodeIdentifier,
	prepFactor hardware.PrepFactor,
	trackCount hardware.TrackCount,
	packName string,
	removable bool,
) error {
	pi := &ioPackets.IoPrepInfo{
		PrepFactor:  prepFactor,
		TrackCount:  trackCount,
		PackName:    packName,
		IsRemovable: removable,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeIdentifier,
		IoFunction:     ioPackets.IofPrep,
		PrepInfo:       pi,
	}
	return io(ch, cp)
}

func ioRead(
	ch *channels.DiskChannel,
	nodeIdentifier hardware.NodeIdentifier,
	buffer []pkg.Word36,
	blockId hardware.BlockId,
	prepFactor hardware.PrepFactor,
) error {
	cw := channels.ControlWord{
		Buffer:    buffer,
		Offset:    0,
		Length:    uint(prepFactor),
		Direction: channels.DirectionForward,
		Format:    channels.TransferPacked,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeIdentifier,
		IoFunction:     ioPackets.IofRead,
		BlockId:        blockId,
		ControlWords:   []channels.ControlWord{cw},
	}
	return io(ch, cp)
}

func ioUnmount(
	ch *channels.DiskChannel,
	nodeIdentifier hardware.NodeIdentifier,
) error {
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeIdentifier,
		IoFunction:     ioPackets.IofUnmount,
	}
	return io(ch, cp)
}

func ioWrite(
	ch *channels.DiskChannel,
	nodeIdentifier hardware.NodeIdentifier,
	buffer []pkg.Word36,
	blockId hardware.BlockId,
	prepFactor hardware.PrepFactor,
) error {
	cw := channels.ControlWord{
		Buffer:    buffer,
		Offset:    0,
		Length:    uint(prepFactor),
		Direction: channels.DirectionForward,
		Format:    channels.TransferPacked,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeIdentifier,
		IoFunction:     ioPackets.IofWrite,
		BlockId:        blockId,
		ControlWords:   []channels.ControlWord{cw},
	}
	return io(ch, cp)
}
