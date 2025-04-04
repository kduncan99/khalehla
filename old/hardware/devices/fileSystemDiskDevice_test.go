// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"os"
	"testing"

	"khalehla/common"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
	ioPackets2 "khalehla/old/hardware/ioPackets"
)

var FileName = "test.pack"

func createPreppedPack(ident uint64, prepFactor uint32, blockSize uint32, blockCount uint64, trackCount uint64) {
	f, _ := os.OpenFile(FileName, os.O_CREATE|os.O_SYNC|os.O_TRUNC|os.O_RDWR, 0666)
	hdr := &fileSystemDiskHeader{
		identifier: ident,
		blockSize:  blockSize,
		prepFactor: prepFactor,
		blockCount: blockCount,
		trackCount: trackCount,
	}

	buffer := make([]byte, 128)
	hdr.serializeInto(buffer)
	_, _ = f.WriteAt(buffer, 0)
	_ = f.Close()
}

func Test_Mount1(t *testing.T) {

	fsd := NewFileSystemDiskDevice(&FileName)
	pkt := &ioPackets2.DiskIoPacket{
		MountInfo: nil,
	}

	fsd.doMount(pkt)
	if pkt.IoStatus != ioPackets2.IosInvalidPacket {
		t.Errorf("Expected invalid packet")
	}
}

func Test_Mount2(t *testing.T) {

	fsd := NewFileSystemDiskDevice(&FileName)
	mi := &ioPackets2.IoMountInfo{
		WriteProtect: false,
	}

	pkt := &ioPackets2.DiskIoPacket{
		IoFunction: ioPackets.IofMount,
		MountInfo:  mi,
	}

	fsd.doMount(pkt)
	if fsd.diskHeader != nil {
		t.Errorf("Expected a nil disk header")
	}

	_ = os.Remove(FileName)
}

func Test_Mount3(t *testing.T) {
	createPreppedPack(fsIdentifier, 28, 128, 640000, 10000)

	fsd := NewFileSystemDiskDevice(&FileName)
	mi := &ioPackets2.IoMountInfo{
		WriteProtect: false,
	}

	pkt := &ioPackets2.DiskIoPacket{
		IoFunction: ioPackets.IofMount,
		MountInfo:  mi,
	}

	fsd.doMount(pkt)
	if fsd.diskHeader == nil {
		t.Errorf("Expected a disk header")
		return
	}

	if fsd.diskHeader.blockSize != 128 {
		t.Errorf("Unexpected block size %v", fsd.diskHeader.blockSize)
	}

	if fsd.diskHeader.prepFactor != 28 {
		t.Errorf("Unexpected prep factor %v", fsd.diskHeader.prepFactor)
	}

	if fsd.diskHeader.blockCount != 640000 {
		t.Errorf("Unexpected block count %v", fsd.diskHeader.blockCount)
	}

	if fsd.diskHeader.trackCount != 10000 {
		t.Errorf("Unexpected track count %v", fsd.diskHeader.trackCount)
	}

	_ = os.Remove(FileName)
}

func Test_IO(t *testing.T) {
	createPreppedPack(fsIdentifier, 28, 128, 640000, 10000)
	fsd := NewFileSystemDiskDevice(&FileName)

	buffer := make([]byte, 128)
	for blockId := hardware.BlockId(0); blockId < 1000; blockId++ {
		common.SerializeUint64IntoBuffer(uint64(blockId), buffer)
		pkt := &ioPackets2.DiskIoPacket{
			Listener:   nil,
			IoFunction: ioPackets.IofWrite,
			IoStatus:   0,
			BlockId:    blockId,
			Buffer:     buffer,
		}

		fsd.StartIo(pkt)
		if pkt.IoStatus != ioPackets2.IosComplete {
			t.Fatalf("Write error:%v", pkt.GetString())
		}
	}

	for blockId := hardware.BlockId(0); blockId < 1000; blockId++ {
		pkt := &ioPackets2.DiskIoPacket{
			Listener:   nil,
			IoFunction: ioPackets.IofRead,
			IoStatus:   0,
			BlockId:    blockId,
			Buffer:     buffer,
		}

		fsd.StartIo(pkt)
		if pkt.IoStatus != ioPackets2.IosComplete {
			t.Fatalf("Read error:%v", pkt.GetString())
		}

		chk := common.DeserializeUint64FromBuffer(buffer)
		if chk != uint64(blockId) {
			t.Fatalf("Block %v chk value was %v", blockId, chk)
		}
	}

	pkt := &ioPackets2.DiskIoPacket{
		IoFunction: ioPackets.IofUnmount,
	}
	fsd.StartIo(pkt)
	if pkt.IoStatus != ioPackets2.IosComplete {
		t.Errorf("Expected IosCompleted:%v", pkt.GetString())
	}

	_ = os.Remove(FileName)
}
