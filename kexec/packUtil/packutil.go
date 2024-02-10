package packUtil

import (
	"errors"
	"fmt"
	"khalehla/kexec/deviceMgr"
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
	if !deviceMgr.IsValidPrepFactor(uint(prepFactor)) {
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

	dd := deviceMgr.NewDiskDevice()
	pkt := deviceMgr.NewDiskIoPacketMount(fileName, false)
	dd.startIo(pkt)
	if pkt.GetIoStatus() != deviceMgr.IosComplete {
		return fmt.Errorf("status %v returned while mounting pack file %v", pkt.GetIoStatus(), fileName)
	}

	pkt = deviceMgr.NewDiskIoPacketPrep(packName, uint(prepFactor), uint(trackCount), removable)
	dd.startIo(pkt)
	if pkt.GetIoStatus() != deviceMgr.IosComplete {
		return fmt.Errorf("status %v returned while prepping pack file %v", pkt.GetIoStatus(), fileName)
	}

	showGeometry(dd)
	showLabelRecord(dd)

	pkt = deviceMgr.NewDiskIoPacketUnmount()
	dd.startIo(pkt)

	return nil
}

func DoShow(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments for show command")
	}

	fileName := args[0]
	if _, err := os.Stat("/path/to/whatever"); err == nil {
		// skip on down
	} else if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file %v does not exist", fileName)
	} else {
		return fmt.Errorf("cannot open file %v:%v", fileName, err)
	}

	dd := deviceMgr.NewDiskDevice()
	pkt := deviceMgr.NewDiskIoPacketMount(fileName, false)
	dd.startIo(pkt)
	if !dd.IsPrepped() {
		return fmt.Errorf("pack is not prepped")
	}

	showGeometry(dd)
	showLabelRecord(dd)

	pkt = deviceMgr.NewDiskIoPacketUnmount()
	dd.startIo(pkt)

	return nil
}

func showLabelRecord(diskDevice *deviceMgr.DiskDevice) {
	label := make([]pkg.Word36, 28)
	pkt := deviceMgr.NewDiskIoPacketReadLabel(label)
	diskDevice.startIo(pkt)
	if pkt.GetIoStatus() == deviceMgr.IosComplete {
		fmt.Println("Label Record:")
		pkg.DumpWord36Buffer(label, 7)
	} else {
		fmt.Printf("Status %v returned while reading label\n", pkt.GetIoStatus())
	}
}

func showGeometry(diskDevice *deviceMgr.DiskDevice) {
	geom := diskDevice.GetGeometry()
	fmt.Println("Geometry:")
	fmt.Printf("  PrepFactor:      %v\n", geom.PrepFactor)
	fmt.Printf("  BlockCount:      %v\n", geom.BlockCount)
	fmt.Printf("  TrackCount:      %v\n", geom.TrackCount)
	fmt.Printf("  BlocksPerTrack:  %v\n", geom.BlocksPerTrack)
	fmt.Printf("  BytesPerBlock:   %v\n", geom.BytesPerBlock)
	fmt.Printf("  SectorsPerBlock: %v\n", geom.SectorsPerBlock)
}
