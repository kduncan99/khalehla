// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tapeUtil

import (
	"fmt"
	"os"

	"khalehla/common"
	pkg2 "khalehla/old/pkg"
)

// BinDevice handles .BIN virtual tape files
//
// A .bin file contains three types of data structures.
// All data structures are, at the outermost, packed 36-bit words.
// That is to say, every 2 logical words is serialized into 9 successive bytes.
// Every data structure is stored on 28-word sector alignments.
// * Each structure begins on a 28-word boundary
// * Each structure is padded if and as necessary to include a whole number of sectors.
//
// The volume header is one sector long. Known fields include:
//   +000 "**TF**" in fieldata, LJSF
//   +001 reel number FD LJSF
//   +002 TDATE$ of creation
//   +003 site name of creator
//   +004
//   +005 sector address of next **EOF*
//   +006 total sectors on tape
//   +007 total files on tape
//   +010 total of what?
//   +011 user-id, or account, or something FD LJSF
//
// The tape block structure is a variable number of sectors long.
// It contains a 6-word header followed by a varying number of data words,
// followed by any necessary padding. The data is 36-bit packed.
// The tape block header is formatted as follows:
//   +000 "**TB**" in fieldata, LJSF
//   +001,H1 file number
//   +001,H2 # sectors of previous block
//   +002 block number starting with 1
//   +003 word length of data portion
//   +004 (unknown)
//   +005 word length of header portion
//
// The tape mark (EOF) structure is one sector long. Fields include:
//   +000 "**EOF*" in fieldata, LJSF
//   +001,H1 file number of file we are ending
//   +001,H2 prev block size
//   +002 number of blocks in previous file
//   +003
//   +004 sector of previous **EOF* or zero if this is the first one
//   +005 sector of next **EOF* or zero if this is the last one

type BinDevice struct {
	file                 *os.File
	currentSectorAddress uint64
}

func NewBinDevice() *BinDevice {
	return &BinDevice{
		file:                 nil,
		currentSectorAddress: 0,
	}
}

func (dev *BinDevice) toByteAddress() int64 {
	wordAddr := dev.currentSectorAddress * 28
	return int64(wordAddr * 9 / 2)
}

func (dev *BinDevice) Close() error {
	if dev.file == nil {
		return fmt.Errorf("file is not open")
	}

	err := dev.file.Close()
	dev.file = nil
	return err
}

func (dev *BinDevice) OpenInputFile(fileName string) error {
	if dev.file != nil {
		return fmt.Errorf("file is already open")
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	dev.file = file
	dev.currentSectorAddress = 0

	return nil
}

func (dev *BinDevice) OpenOutputFile(fileName string) error {
	if dev.file != nil {
		return fmt.Errorf("file is already open")
	}

	return nil // TODO
}

func (dev *BinDevice) ReadVolumeHeader() (volHeader *VolumeHeader, err error) {
	volHeader = nil
	err = nil
	if dev.file == nil {
		err = fmt.Errorf("file is not open")
		return
	}

	words := make([]pkg2.Word36, 28)
	bytes := make([]byte, 126)
	offset := dev.toByteAddress()
	_, err = dev.file.ReadAt(bytes, offset)
	if err != nil {
		return nil, err
	}

	common.UnpackWord36Strict(bytes, words)
	sentinel := words[0].ToStringAsFieldata()
	if sentinel != "**TF**" {
		err = fmt.Errorf("found '%v' - expected '**TF**'", sentinel)
	}

	dev.currentSectorAddress++
	volHeader = NewVolumeHeader().
		SetReelNumber(words[01].ToStringAsFieldata()).
		SetSiteId(words[03].ToStringAsFieldata()).
		SetAccountId(words[011].ToStringAsFieldata()).
		SetNumberOfFiles(words[07].GetW())

	return
}

func (dev *BinDevice) ReadFileHeader() (fileHeader *FileHeader, err error) {
	fileHeader = nil
	err = nil
	if dev.file == nil {
		err = fmt.Errorf("file is not open")
		return
	}

	// This format does not store individual file headers,
	// therefore we cannot read a header, and must produce a default FileHeader struct.
	fileHeader = NewFileHeader()
	return
}

func (dev *BinDevice) ReadBlock() (data []pkg2.Word36, eof bool, err error) {
	data = nil
	eof = false
	err = nil
	if dev.file == nil {
		err = fmt.Errorf("file is not open")
		return
	}

	// read 6-word block header first, so we know how big the entire block is
	header := make([]pkg2.Word36, 6)
	bytes := make([]byte, 6*9/2)
	offset := dev.toByteAddress()
	_, err = dev.file.ReadAt(bytes, offset)
	if err != nil {
		return
	}

	common.UnpackWord36Strict(bytes, header)
	sentinel := header[0].ToStringAsFieldata()
	if sentinel == "**EOF*" {
		eof = true
		dev.currentSectorAddress++
		return
	} else if sentinel != "**TB**" {
		err = fmt.Errorf("found '%v' - expected '**TB**'", sentinel)
	}

	// read the entire block
	headerWords := header[05].GetW()
	dataWords := header[03].GetW()
	totalWords := headerWords + dataWords
	mod := totalWords % 28
	if mod != 0 {
		totalWords += 28 - mod
	}

	dataBlock := make([]pkg2.Word36, totalWords)
	bytes = make([]byte, len(dataBlock)*9/2)
	_, err = dev.file.ReadAt(bytes, offset)
	if err != nil {
		return
	}
	dev.currentSectorAddress += totalWords / 28

	common.UnpackWord36Strict(bytes, dataBlock)
	data = dataBlock[headerWords : headerWords+dataWords]
	return
}

func (dev *BinDevice) WriteVolumeHeader() error {
	if dev.file == nil {
		return fmt.Errorf("file is not open")
	}

	// TODO

	return nil
}

func (dev *BinDevice) WriteFileHeader() error {
	if dev.file == nil {
		return fmt.Errorf("file is not open")
	}

	// TODO

	return nil
}

func (dev *BinDevice) WriteBlock(buffer []pkg2.Word36) error {
	if dev.file == nil {
		return fmt.Errorf("file is not open")
	}

	// TODO

	return nil
}

func (dev *BinDevice) WriteEndOfFile() error {
	if dev.file == nil {
		return fmt.Errorf("file is not open")
	}

	return nil
}
