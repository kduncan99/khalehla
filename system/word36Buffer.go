package system

import "fmt"

type Word36Buffer struct {
	Value []uint64
}

func (buffer *Word36Buffer) packInto(sourceOffset int, length int, dest []byte, destOffset int) error {
	//	sourceOffset is in words
	//	length is in words, and must be a multiple of two
	//	destOffset is in bytes

	if length&0x01 != 0 {
		return fmt.Errorf("length is not an even number")
	}

	//TODO verify we do not exceed any array boundaries

	sx := sourceOffset
	dx := destOffset
	for remaining := length; remaining > 0; remaining -= 2 {
		w1 := buffer.Value[sx]
		w2 := buffer.Value[sx+1]
		sx += 2

		b8 := (byte)(w2 & 0xff)
		w2 >>= 8
		b7 := (byte)(w2 & 0xff)
		w2 >>= 8
		b6 := (byte)(w2 & 0xff)
		w2 >>= 8
		b5 := (byte)(w2 & 0xff)
		w2 >>= 8
		b4 := (byte)(w2&0xf)<<4 | (byte)(w1&0x0f)
		w1 >>= 4
		b3 := (byte)(w1 & 0xff)
		w1 >>= 8
		b2 := (byte)(w1 & 0xff)
		w1 >>= 8
		b1 := (byte)(w1 & 0xff)
		w1 >>= 8
		b0 := (byte)(w1 & 0xff)

		dest[dx] = b0
		dx++
		dest[dx] = b1
		dx++
		dest[dx] = b2
		dx++
		dest[dx] = b3
		dx++
		dest[dx] = b4
		dx++
		dest[dx] = b5
		dx++
		dest[dx] = b6
		dx++
		dest[dx] = b7
		dx++
		dest[dx] = b8
		dx++
	}

	return nil
}

func (buffer *Word36Buffer) unpackFrom(source []byte, sourceOffset int, length int, destOffset int) error {
	//	destOffset is in words
	//	length is in words, and must be a multiple of two
	//	sourceOffset is in bytes

	if length&0x01 != 0 {
		return fmt.Errorf("length is not an even number")
	}

	//TODO verify we do not exceed any array boundaries

	sx := sourceOffset
	dx := destOffset
	for remaining := length; remaining > 0; remaining -= 2 {
		w1 := (uint64)(source[sx])
		w1 <<= 8
		sx++
		w1 |= (uint64)(source[sx])
		w1 <<= 8
		sx++
		w1 |= (uint64)(source[sx])
		w1 <<= 8
		sx++
		w1 |= (uint64)(source[sx])
		w1 <<= 8
		sx++
		w1 |= (uint64)(source[sx] >> 4)

		w2 := (uint64)(source[sx] & 0xF)
		w2 <<= 4
		sx++
		w2 |= (uint64)(source[sx])
		w2 <<= 8
		sx++
		w2 |= (uint64)(source[sx])
		w2 <<= 8
		sx++
		w2 |= (uint64)(source[sx])
		w2 <<= 8
		sx++
		w2 |= (uint64)(source[sx])
		sx++

		buffer.Value[dx] = w1
		buffer.Value[dx+1] = w2
		dx += 2
	}

	return nil
}
