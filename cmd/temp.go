package main

import (
	"fmt"
	"unsafe"
)

type Word uint64

type Packet struct {
	w0 Word
	w1 Word
	w2 Word
}

// _, err := bd.file.ReadAt(unsafe.Slice((*byte)(unsafe.Pointer(&buffer[0])), bd.geometry.bytesPerBlock), pos)
func writePacketToBuffer(packet Packet, buffer []Word) {
	pSource := unsafe.Slice((*Word)(unsafe.Pointer(&packet)), 3)
	// pDest := unsafe.Slice((*Word)(unsafe.Pointer(&buffer[0])), len(buffer))
	for wx := 0; wx < 3; wx++ {
		buffer[wx] = pSource[wx]
	}
}

func main() {
	p := Packet{10, 20, 30}
	b := make([]Word, 10)
	writePacketToBuffer(p, b)
	fmt.Printf("%v\n", b)
}
