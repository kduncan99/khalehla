// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package packUtil

import (
	"errors"
	"fmt"
	"khalehla/kexec/deviceMgr"
	"khalehla/kexec/types"
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
	if !deviceMgr.IsValidPackName(packName) {
		return fmt.Errorf("invalid pack name")
	}

	prepFactor, err := strconv.Atoi(args[2])
	if err != nil || prepFactor <= 0 {
		return fmt.Errorf("error in prepfactor argument")
	}
	if !deviceMgr.IsValidPrepFactor(types.PrepFactor(prepFactor)) {
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

	dc := deviceMgr.NewDiskChannel()
	dd := deviceMgr.NewDiskDevice(nil)
	ni := types.NodeIdentifier(pkg.NewFromStringToFieldata("DISK0", 1)[0])
	_ = dc.AssignDevice(ni, dd)

	pkt := deviceMgr.NewDiskIoPacketMount(ni, fileName, false)
	dc.StartIo(pkt)
	if pkt.GetIoStatus() != types.IosComplete {
		return fmt.Errorf("status %v returned while mounting pack file %v", pkt.GetIoStatus(), fileName)
	}

	pkt = deviceMgr.NewDiskIoPacketPrep(ni, packName, types.PrepFactor(prepFactor), types.TrackCount(trackCount), removable)
	dc.StartIo(pkt)
	if pkt.GetIoStatus() != types.IosComplete {
		return fmt.Errorf("status %v returned while prepping pack file %v", pkt.GetIoStatus(), fileName)
	}

	showLabelRecord(dc, ni, true)

	pkt = deviceMgr.NewDiskIoPacketUnmount(ni)
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

	dc := deviceMgr.NewDiskChannel()
	dd := deviceMgr.NewDiskDevice(nil)
	ni := types.NodeIdentifier(pkg.NewFromStringToFieldata("DISK0", 1)[0])
	_ = dc.AssignDevice(ni, dd)

	pkt := deviceMgr.NewDiskIoPacketMount(ni, fileName, false)
	dc.StartIo(pkt)
	if !dd.IsPrepped() {
		return fmt.Errorf("pack is not prepped")
	}

	showLabelRecord(dc, ni, true)

	pkt = deviceMgr.NewDiskIoPacketUnmount(ni)
	dc.StartIo(pkt)

	return nil
}

func showLabelRecord(channel types.Channel, nodeId types.NodeIdentifier, interpret bool) {
	label := make([]pkg.Word36, 28)
	pkt := deviceMgr.NewDiskIoPacketReadLabel(nodeId, label)
	channel.StartIo(pkt)
	if pkt.GetIoStatus() != types.IosComplete {
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
	fmt.Printf("S0+S1+HMBT+Pad:       %d words\n", label[011].GetH1())
	fmt.Printf("Master Bit Table Len: %d words\n", label[011].GetH2())
	fmt.Printf("VOL1 Version:         %d\n", label[014].GetS2())
	fmt.Printf("Disk Capacity:        %d tracks\n", label[016].GetW())
	fmt.Printf("Words Per Phys Record:%d\n", label[017].GetH1())
	fmt.Printf("Total Blocks:         %d\n", label[021].GetW())
}
