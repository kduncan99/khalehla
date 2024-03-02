// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package packUtil

import (
	"errors"
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
	"os"
	"strconv"
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
	if !kexec.IsValidPackName(packName) {
		return fmt.Errorf("invalid pack name")
	}

	prepFactor, err := strconv.Atoi(args[2])
	if err != nil || prepFactor <= 0 {
		return fmt.Errorf("error in prepfactor argument")
	}
	if !kexec.IsValidPrepFactor(kexec.PrepFactor(prepFactor)) {
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

	dc := nodeMgr.NewDiskChannel()
	dd := nodeMgr.NewFileSystemDiskDevice(nil)
	devId := kexec.NodeIdentifier(pkg.NewFromStringToFieldata("DISK0", 1)[0])
	_ = dc.AssignDevice(devId, dd)

	pkt := nodeMgr.NewDiskIoPacketMount(devId, fileName, false)
	dc.StartIo(pkt)
	if pkt.GetIoStatus() != nodeMgr.IosComplete {
		return fmt.Errorf("status %v returned while mounting pack file %v", pkt.GetIoStatus(), fileName)
	}

	pkt = nodeMgr.NewDiskIoPacketPrep(devId, packName, kexec.PrepFactor(prepFactor), kexec.TrackCount(trackCount), removable)
	dc.StartIo(pkt)
	if pkt.GetIoStatus() != nodeMgr.IosComplete {
		return fmt.Errorf("status %v returned while prepping pack file %v", pkt.GetIoStatus(), fileName)
	}

	showLabelRecord(dc, devId, true)

	pkt = nodeMgr.NewDiskIoPacketUnmount(devId)
	dc.StartIo(pkt)

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

	dc := nodeMgr.NewDiskChannel()
	dd := nodeMgr.NewFileSystemDiskDevice(nil)
	devId := kexec.NodeIdentifier(pkg.NewFromStringToFieldata("DISK0", 1)[0])
	_ = dc.AssignDevice(devId, dd)

	pkt := nodeMgr.NewDiskIoPacketMount(devId, fileName, false)
	dc.StartIo(pkt)
	if !dd.IsPrepped() {
		return fmt.Errorf("pack is not prepped")
	}

	showLabelRecord(dc, devId, true)

	pkt = nodeMgr.NewDiskIoPacketUnmount(devId)
	dc.StartIo(pkt)

	return nil
}

func showLabelRecord(channel nodeMgr.Channel, devId kexec.NodeIdentifier, interpret bool) {
	label := make([]pkg.Word36, 28)
	pkt := nodeMgr.NewDiskIoPacketReadLabel(devId, label)
	channel.StartIo(pkt)
	if pkt.GetIoStatus() != nodeMgr.IosComplete {
		fmt.Printf("Status %v returned while reading label\n", pkt.GetIoStatus())
		return
	}

	fmt.Println("Label Record:")
	pkg.DumpWord36Buffer(label, 7)
	if !interpret {
		return
	}

	fmt.Printf("Pack Name:            %v%v\n", label[1].ToStringAsAscii(), label[2].ToStringAsAscii())
	fmt.Printf("First Dir Track DRWA: %v\n", label[3].ToStringAsOctal())
	fmt.Printf("Records Per Track:    %d\n", label[4].GetH1())
	fmt.Printf("Words Per Record:     %d\n", label[4].GetH2())
	fmt.Printf("VOL1 Version:         %d\n", label[014].GetS2())
	fmt.Printf("Disk Capacity:        %d tracks\n", label[016].GetW())
	fmt.Printf("Words Per Phys Record:%d\n", label[017].GetH1())
	fmt.Printf("Total Blocks:         %d\n", label[021].GetW())
}
