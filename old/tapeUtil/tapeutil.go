// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tapeUtil

import (
	"fmt"
)

// kexec tape file format (all offsets and lengths are in bytes, offsets are from the start of the file):
// tapes are composed of data blocks and tape marks.
// A data block has a header and a payload portion, and is formatted as such:
//
//	+0 32-bit length of payload 0 < n < 0xFFFFFFFE
//	+4 32-bit length of payload of previous block, or zero if this is the first block on the tape
//	       or following a tape mark
//	+8 payload starts here
//
// A tape mark has only a header portion, and is formatted as such:
//
//	+0 0xFFFFFFFF
//	+4 32-bit length of payload of previous block, or zero if this is the first thing on the tape,
//	       or this tape mark follows another tape mark
// type tapeFileInfo struct {
// 	reelNumber  string
// 	owner       string
// 	siteName    string
// 	sectorCount uint
// 	fileCount   uint
// }

// func makeDataBlockControlWord(payloadLength uint64, previousLength uint64) []byte {
// 	cw := make([]byte, 8)
// 	cw[0] = byte(payloadLength >> 24)
// 	cw[1] = byte(payloadLength >> 16)
// 	cw[2] = byte(payloadLength >> 8)
// 	cw[3] = byte(payloadLength)
// 	cw[4] = byte(previousLength >> 24)
// 	cw[5] = byte(previousLength >> 16)
// 	cw[6] = byte(previousLength >> 8)
// 	cw[7] = byte(previousLength)
// 	return cw
// }

// func makeTapeMarkControlWord(previousLength uint64) []byte {
// 	cw := make([]byte, 8)
// 	cw[0] = 0xFF
// 	cw[1] = 0xFF
// 	cw[2] = 0xFF
// 	cw[3] = 0xFF
// 	cw[4] = byte(previousLength >> 24)
// 	cw[5] = byte(previousLength >> 16)
// 	cw[6] = byte(previousLength >> 8)
// 	cw[7] = byte(previousLength)
// 	return cw
// }

// func copyFile(fIn *os.File, startingWordOffset uint64, fOut *os.File) (finalWordOffset uint64, err error) {
// 	previousLength := uint64(0)
// 	wordOffset := startingWordOffset
// 	blockCount := uint64(0)
// 	wordCount := uint64(0)
// 	for {
// 		block, err := readVariableBlock(fIn, wordOffset, 6)
// 		if err != nil {
// 			return 0, err
// 		}
//
// 		sentinel := block[0].ToStringAsFieldata()
// 		if sentinel == "**EOF*" {
// 			wordOffset += 28
// 			break
// 		} else if sentinel != "**TB**" {
// 			return 0, fmt.Errorf("unrecognized block: %v", sentinel)
// 		}
//
// 		headerWordCount := block[5].GetW()
// 		dataWordCount := block[3].GetW()
// 		totalWordCount := headerWordCount + dataWordCount
// 		mod := totalWordCount % 28
// 		if mod != 0 {
// 			totalWordCount += 28 - mod
// 		}
//
// 		block, err = readVariableBlock(fIn, startingWordOffset, int(totalWordCount))
// 		if err != nil {
// 			return 0, err
// 		}
//
// 		dataByteCount := dataWordCount * 9 / 2
// 		_, err = fOut.Write(makeDataBlockControlWord(dataByteCount, previousLength))
// 		if err != nil {
// 			return 0, err
// 		}
//
// 		byteBuffer := make([]byte, dataByteCount)
// 		pkg.PackWord36(block[headerWordCount:headerWordCount+dataWordCount], byteBuffer)
// 		_, err = fOut.Write(byteBuffer)
// 		if err != nil {
// 			return 0, err
// 		}
//
// 		wordOffset += totalWordCount
// 		blockCount++
// 		wordCount += dataWordCount
// 	}
//
// 	_, err = fOut.Write(makeTapeMarkControlWord(previousLength))
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	fmt.Printf("Copy file complete: %v blocks %v words\n", blockCount, wordCount)
// 	return wordOffset, nil
// }

// func readTapeHeader(file *os.File) (tfInfo *tapeFileInfo, wordLength uint64, err error) {
// 	header, err := readVariableBlock(file, 0, 28)
// 	if err != nil {
// 		return nil, 0, err
// 	}
//
// 	tfInfo = &tapeFileInfo{
// 		reelNumber:  header[1].ToStringAsFieldata(),
// 		siteName:    header[3].ToStringAsFieldata(),
// 		sectorCount: uint(header[6].GetW()),
// 		fileCount:   uint(header[7].GetW()),
// 		owner:       header[011].ToStringAsFieldata(),
// 	}
//
// 	return tfInfo, 28, nil
// }

// func readVariableBlock(file *os.File, wordOffset uint64, size int) ([]pkg.Word36, error) {
// 	buffer := make([]pkg.Word36, size)
// 	bytes := make([]byte, len(buffer)*9/2)
// 	_, err := file.ReadAt(bytes, int64(wordOffset*9/2))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	pkg.UnpackWord36(bytes, buffer)
// 	return buffer, nil
// }

// func analyze(f *os.File) {
// 	blockAddress := int64(0)
// 	sectorAddress := int64(0)
// 	done := false
// 	for !done {
// 		sector, err := readVariableBlock(f, blockAddress, 28)
// 		if err != nil {
// 			fmt.Printf("Error:%v\n", err)
// 			break
// 		}
//
// 		fmt.Printf("Sector %012o (%012o)\n", sectorAddress, blockAddress)
// 		switch sector[0].ToStringAsFieldata() {
// 		case "**TF**":
// 			pkg.DumpWord36Buffer(sector, 7)
// 			blockAddress += 28
// 			sectorAddress++
//
// 		case "**TB**":
// 			pkg.DumpWord36Buffer(sector, 7)
// 			headerWords := int64(sector[5].GetW())
// 			dataWords := int64(sector[3].GetW())
// 			totalWords := headerWords + dataWords
// 			for totalWords%28 != 0 {
// 				totalWords++
// 			}
// 			sectorAddress += totalWords / 28
// 			blockAddress += totalWords
//
// 		case "**EOF*":
// 			pkg.DumpWord36Buffer(sector, 7)
// 			blockAddress += 28
// 			sectorAddress++
//
// 		default:
// 			fmt.Println("Unknown block format:")
// 			pkg.DumpWord36Buffer(sector, 7)
// 			done = true
// 		}
// 	}
// }

func DoArguments() error {
	// TODO
	return nil
}

func DoUsage() {
	fmt.Println("Usage: tapeUtil {switches} {command} {file-name-1} [ {file-name-2} ]")
	fmt.Println("  tapeConvert {input_file} {output_file} [ -if {format} ] [ -of {format} ]")
	fmt.Println("-if format - defines the format of the input file")
	fmt.Println("  format: BIN | KEXEC | KFAST ")
	fmt.Println("-of format - defines the format of the output file")
	fmt.Println("  format: BIN | KEXEC | KFAST ")
	fmt.Println("-il - indicates that the input file is a labeled tape image")
	fmt.Println("-ol - indicates that the output file should be a labeled tape image")
}

func DoMain(args []string) {
	// fIn, err := os.Open("/Users/kduncan/Desktop/VT_1234.bin")
	// if err != nil {
	// 	panic(err)
	// }
	// defer fIn.Close()
	//
	// fOut, err := os.OpenFile("tapeVolume", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// defer fOut.Close()
	//
	// tfInfo, wordOffset, err := readTapeHeader(fIn)
	// if err != nil {
	// 	panic(err)
	// }

	inputDev := NewBinDevice()
	defer inputDev.Close()
	err := inputDev.OpenInputFile("/Users/kduncan/Desktop/VT_1234.bin")
	if err != nil {
		panic(err)
	}

	outputDev := NewFSTapeDevice()
	defer outputDev.Close()
	err = outputDev.OpenOutputFile("001234.tape")
	if err != nil {
		panic(err)
	}

	vh, err := inputDev.ReadVolumeHeader()
	if err != nil {
		panic(err)
	}

	fmt.Printf("**Tape Header Info:\n")
	fmt.Printf("  Reel Number:  %v\n", vh.reelNumber)
	fmt.Printf("  Owner:        %v\n", vh.accountId)
	fmt.Printf("  Site Name:    %v\n", vh.siteId)
	fmt.Printf("  File Count:   %v\n", vh.numberOfFiles)

	endOfVolume := false
	for !endOfVolume {
		_, err := inputDev.ReadFileHeader()
		if err != nil {
			fmt.Printf("Error:%v\n", err)
			endOfVolume = true
			goto done
		}

		for {
			data, eof, err := inputDev.ReadBlock()
			if err != nil {
				fmt.Printf("Error:%v\n", err)
				endOfVolume = true
				goto done
			}

			if eof {
				fmt.Printf("**EOF\n")
				err = outputDev.WriteEndOfFile()
				if err != nil {
					fmt.Printf("Error:%v\n", err)
					endOfVolume = true
					goto done
				}
				break
			} else {
				fmt.Printf("**Data: %v words\n", len(data))
				err = outputDev.WriteBlock(data)
				if err != nil {
					fmt.Printf("Error:%v\n", err)
					endOfVolume = true
					goto done
				}
			}
		}
	}

done:
}
